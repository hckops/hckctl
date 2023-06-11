package docker

import (
	"context"
	"io"
	"strings"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
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
	eventBus.Publish(newInitClientDockerEvent())

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
	cli.eventBus.Publish(newCloseClientDockerEvent())
	cli.eventBus.Close()
	return cli.docker.Close()
}

type SetupImageOpts struct {
	ImageName           string
	OnPullImageCallback func()
}

func (cli *DockerClient) Setup(opts *SetupImageOpts) error {
	cli.eventBus.Publish(newSetupImageDockerEvent(opts.ImageName))

	// TODO delete dangling images

	reader, err := cli.docker.ImagePull(cli.ctx, opts.ImageName, types.ImagePullOptions{})
	if err != nil {
		return errors.Wrap(err, "error image pull")
	}
	defer reader.Close()

	cli.eventBus.Publish(newPullImageDockerEvent(opts.ImageName))
	opts.OnPullImageCallback()

	// suppress default output
	if _, err := io.Copy(io.Discard, reader); err != nil {
		return errors.Wrap(err, "error image pull output message")
	}

	return nil
}

type CreateContainerOpts struct {
	ContainerName   string
	ContainerConfig *container.Config
	HostConfig      *container.HostConfig
}

func (cli *DockerClient) CreateContainer(opts *CreateContainerOpts) (string, error) {
	cli.eventBus.Publish(newCreateContainerDockerEvent(opts.ContainerName))

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
}

func (cli *DockerClient) ExecContainer(opts *ExecContainerOpts) error {
	cli.eventBus.Publish(newExecContainerDockerEvent(opts.ContainerId))

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
	if err != nil {
		return errors.Wrap(err, "error container exec attach")
	}
	defer execAttachResponse.Close()

	removeContainerCallback := func() {
		if err := cli.RemoveContainer(opts.ContainerId); err != nil {
			cli.eventBus.Publish(newExecContainerErrorDockerEvent(opts.ContainerId, errors.Wrap(err, "remove container")))
		}
	}

	onStreamErrorCallback := func(err error) {
		cli.eventBus.Publish(newExecContainerErrorDockerEvent(opts.ContainerId, errors.Wrap(err, "stream container")))
	}

	handleStreams(opts, &execAttachResponse, removeContainerCallback, onStreamErrorCallback)

	// fixes echoes and handle SIGTERM interrupt properly
	if terminal, err := util.NewRawTerminal(opts.InStream); err == nil {
		defer terminal.Restore()
	}

	cli.eventBus.Publish(newExecContainerWaitingDockerEvent(opts.ContainerId))
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
	cli.eventBus.Publish(newRemoveContainerDockerEvent(containerId))

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

// TODO filter with prefix

type DockerContainerInfo struct {
	ContainerId   string
	ContainerName string
	// TODO command, ports, check status, etc
}

func (cli *DockerClient) ListContainers() ([]DockerContainerInfo, error) {

	containers, err := cli.docker.ContainerList(cli.ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "error docker list")
	}

	var result []DockerContainerInfo
	for index, c := range containers {
		cli.eventBus.Publish(newListContainersDockerEvent(index, c.ID, c.Names[0]))
		result = append(result, DockerContainerInfo{ContainerId: c.ID, ContainerName: c.Names[0]})
	}

	return result, nil
}
