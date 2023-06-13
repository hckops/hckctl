package docker

import (
	"context"
	"io"
	"strings"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	dockerApi "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/client"
	"github.com/hckops/hckctl/pkg/util"
)

type DockerClient struct {
	ctx      context.Context
	docker   *dockerApi.Client
	eventBus *client.EventBus
}

func NewDockerClient(eventBus *client.EventBus) (*DockerClient, error) {
	eventBus.Publish(newClientInitDockerEvent())

	dockerClient, err := dockerApi.NewClientWithOpts(dockerApi.FromEnv, dockerApi.WithAPIVersionNegotiation())
	if err != nil {
		return nil, errors.Wrap(err, "error docker client")
	}

	return &DockerClient{
		ctx:      context.Background(),
		docker:   dockerClient,
		eventBus: eventBus,
	}, nil
}

func (cli *DockerClient) Close() error {
	cli.eventBus.Publish(newClientCloseDockerEvent())
	cli.eventBus.Close()
	return cli.docker.Close()
}

type SetupImageOpts struct {
	ImageName           string
	OnPullImageCallback func()
}

func (cli *DockerClient) Setup(opts *SetupImageOpts) error {
	cli.eventBus.Publish(newImageSetupDockerEvent(opts.ImageName))

	reader, err := cli.docker.ImagePull(cli.ctx, opts.ImageName, types.ImagePullOptions{})
	if err != nil {
		return errors.Wrap(err, "error image pull")
	}
	defer reader.Close()

	cli.eventBus.Publish(newImagePullDockerEvent(opts.ImageName))
	opts.OnPullImageCallback()

	// suppress default output
	if _, err := io.Copy(io.Discard, reader); err != nil {
		return errors.Wrap(err, "error image pull output")
	}

	// cleanup old images
	if err := cli.removeDanglingImages(); err != nil {
		return errors.Wrap(err, "error setup cleanup")
	}

	return nil
}

func (cli *DockerClient) removeDanglingImages() error {

	// dangling images have no tags <none>
	images, err := cli.docker.ImageList(cli.ctx, types.ImageListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key: "dangling", Value: "true",
		}),
	})
	if err != nil {
		return errors.Wrap(err, "error image list dangling")
	}
	for _, image := range images {
		cli.eventBus.Publish(newImageRemoveDockerEvent(image.ID))

		// do not force: there might be running containers with old images
		_, err := cli.docker.ImageRemove(cli.ctx, image.ID, types.ImageRemoveOptions{})
		if err != nil {
			// ignore failures
			cli.eventBus.Publish(newImageRemoveErrorDockerEvent(image.ID, err))
		}
	}

	return nil
}

type CreateContainerOpts struct {
	ContainerName   string
	ContainerConfig *container.Config
	HostConfig      *container.HostConfig
}

func (cli *DockerClient) CreateContainer(opts *CreateContainerOpts) (string, error) {
	cli.eventBus.Publish(newContainerCreateDockerEvent(opts.ContainerName))

	newContainer, err := cli.docker.ContainerCreate(
		cli.ctx,
		opts.ContainerConfig,
		opts.HostConfig,
		nil, // networkingConfig
		nil, // platform
		opts.ContainerName)
	if err != nil {
		return "", errors.Wrap(err, "error container create")
	}

	if err := cli.docker.ContainerStart(cli.ctx, newContainer.ID, types.ContainerStartOptions{}); err != nil {
		return "", errors.Wrap(err, "error container start")
	}

	return newContainer.ID, nil
}

type ExecContainerOpts struct {
	ContainerId                string
	Shell                      string
	InStream                   io.Reader
	OutStream                  io.Writer
	ErrStream                  io.Writer
	IsTty                      bool
	OnContainerWaitingCallback func()
	OnExitCallback             func()
}

// TODO handle distroless i.e. shell == none
// TODO issue with powershell i.e. /usr/bin/pwsh

func (cli *DockerClient) ExecContainer(opts *ExecContainerOpts) error {
	cli.eventBus.Publish(newContainerExecDockerEvent(opts.ContainerId))

	// default shell
	var shellCmd string
	if strings.TrimSpace(opts.Shell) != "" {
		shellCmd = opts.Shell
	} else {
		shellCmd = "/bin/bash"
	}

	execCreateResponse, err := cli.docker.ContainerExecCreate(cli.ctx, opts.ContainerId, types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Detach:       false,
		Tty:          opts.IsTty,
		Cmd:          []string{shellCmd},
	})
	if err != nil {
		return errors.Wrap(err, "error container exec create")
	}

	execAttachResponse, err := cli.docker.ContainerExecAttach(cli.ctx, execCreateResponse.ID, types.ExecStartCheck{
		Tty: opts.IsTty,
	})
	// TODO command, ports, check status, etc
	if err != nil {
		return errors.Wrap(err, "error container exec attach")
	}
	defer execAttachResponse.Close()

	onStreamErrorCallback := func(err error) {
		cli.eventBus.Publish(newContainerExecErrorDockerEvent(opts.ContainerId, errors.Wrap(err, "stream container")))
	}
	// TODO move back internally OnExitCallback and refactor ExecWait/ExecWaitRemove vs Exec: issue WaitConditionNotRunning
	// newExecContainerExitDockerEvent
	if opts.OnExitCallback == nil {
		opts.OnExitCallback = func() {
			cli.docker.ContainerRestart(cli.ctx, opts.ContainerId, container.StopOptions{})
		}
	}

	handleStreams(opts, &execAttachResponse, opts.OnExitCallback, onStreamErrorCallback)

	// fixes echoes and handle SIGTERM interrupt properly
	if terminal, err := util.NewRawTerminal(opts.InStream); err == nil {
		defer terminal.Restore()
	}

	cli.eventBus.Publish(newContainerExecWaitDockerEvent(opts.ContainerId))
	opts.OnContainerWaitingCallback()

	// waits for interrupt signals
	statusCh, errCh := cli.docker.ContainerWait(cli.ctx, opts.ContainerId, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return errors.Wrap(err, "error container wait")
		}
	case <-statusCh:
	}
	return nil
}

func (cli *DockerClient) RemoveContainer(containerId string) error {
	cli.eventBus.Publish(newContainerRemoveDockerEvent(containerId))

	if err := cli.docker.ContainerRemove(cli.ctx, containerId, types.ContainerRemoveOptions{Force: true}); err != nil {
		return errors.Wrap(err, "error docker remove")
	}
	return nil
}

func handleStreams(
	opts *ExecContainerOpts,
	execAttachResponse *types.HijackedResponse,
	onCloseCallback func(),
	onStreamErrorCallback func(error),
) {

	var once sync.Once
	go func() {

		if opts.IsTty {
			if _, err := io.Copy(opts.OutStream, execAttachResponse.Reader); err != nil {
				onStreamErrorCallback(errors.Wrap(err, "error copy stdout docker->local"))
			}
		} else {
			if _, err := stdcopy.StdCopy(opts.OutStream, opts.ErrStream, execAttachResponse.Reader); err != nil {
				onStreamErrorCallback(errors.Wrap(err, "error copy stdout and stderr docker->local"))
			}
		}

		once.Do(onCloseCallback)
	}()
	go func() {
		if _, err := io.Copy(execAttachResponse.Conn, opts.InStream); err != nil {
			onStreamErrorCallback(errors.Wrap(err, "error copy stdin local->docker"))
		}

		once.Do(onCloseCallback)
	}()
}

type DockerContainerInfo struct {
	ContainerId   string
	ContainerName string
}

func (cli *DockerClient) ListContainers(namePrefix string) ([]DockerContainerInfo, error) {

	containers, err := cli.docker.ContainerList(cli.ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key: "name", Value: namePrefix,
		}),
	})
	if err != nil {
		return nil, errors.Wrap(err, "error container list")
	}

	var result []DockerContainerInfo
	for index, c := range containers {
		cli.eventBus.Publish(newContainerListDockerEvent(index, c.ID, c.Names[0]))
		result = append(result, DockerContainerInfo{ContainerId: c.ID, ContainerName: c.Names[0]})
	}

	return result, nil
}
