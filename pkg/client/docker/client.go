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

	"github.com/hckops/hckctl/pkg/util"
)

func NewDockerClient() (*DockerClient, error) {

	dockerClient, err := dockerApi.NewClientWithOpts(dockerApi.FromEnv, dockerApi.WithAPIVersionNegotiation())
	if err != nil {
		return nil, errors.Wrap(err, "error docker client")
	}

	return &DockerClient{
		ctx:    context.Background(),
		docker: dockerClient,
	}, nil
}

func (client *DockerClient) Close() error {
	return client.docker.Close()
}

func (client *DockerClient) ImagePull(opts *ImagePullOpts) error {

	reader, err := client.docker.ImagePull(client.ctx, opts.ImageName, types.ImagePullOptions{})
	if err != nil {
		return errors.Wrap(err, "error image pull")
	}
	defer reader.Close()

	opts.OnImagePullCallback()

	// suppress default output
	if _, err := io.Copy(io.Discard, reader); err != nil {
		return errors.Wrap(err, "error image pull output")
	}

	return nil
}

func (client *DockerClient) ImageRemoveDangling(opts *ImageRemoveOpts) error {

	// dangling images have no tags <none>
	images, err := client.docker.ImageList(client.ctx, types.ImageListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key: "dangling", Value: "true",
		}),
	})
	if err != nil {
		return errors.Wrap(err, "error image list dangling")
	}
	for _, image := range images {
		opts.OnImageRemoveCallback(image.ID)

		_, err := client.docker.ImageRemove(client.ctx, image.ID, types.ImageRemoveOptions{})
		if err != nil {
			// ignore failures: there might be running containers with old images
			opts.OnImageRemoveErrorCallback(image.ID, err)
		}
	}

	return nil
}

func (client *DockerClient) ContainerCreate(opts *ContainerCreateOpts) (string, error) {

	newContainer, err := client.docker.ContainerCreate(
		client.ctx,
		opts.ContainerConfig,
		opts.HostConfig,
		nil, // networkingConfig
		nil, // platform
		opts.ContainerName)
	if err != nil {
		return "", errors.Wrap(err, "error container create")
	}

	if err := client.docker.ContainerStart(client.ctx, newContainer.ID, types.ContainerStartOptions{}); err != nil {
		return "", errors.Wrap(err, "error container start")
	}

	return newContainer.ID, nil
}

// TODO handle distroless i.e. shell == none
// TODO issue with powershell i.e. /usr/bin/pwsh
func defaultShell(command string) string {
	if shellCmd := strings.TrimSpace(command); shellCmd != "" {
		return shellCmd
	} else {
		return "/bin/bash"
	}
}

func (client *DockerClient) ContainerExec(opts *ContainerExecOpts) error {

	// TODO ContainerExecKill https://github.com/moby/moby/pull/41548
	// TODO detach streams properly https://github.com/docker/cli/blob/master/cli/command/container/exec.go

	return nil
}

func (client *DockerClient) ContainerAttach(opts *ContainerAttachOpts) error {

	execCreateResponse, err := client.docker.ContainerExecCreate(client.ctx, opts.ContainerId, types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Detach:       false,
		Tty:          opts.IsTty,
		Cmd:          []string{defaultShell(opts.Shell)},
	})
	if err != nil {
		return errors.Wrap(err, "error container exec create")
	}

	execAttachResponse, err := client.docker.ContainerExecAttach(client.ctx, execCreateResponse.ID, types.ExecStartCheck{
		Tty: opts.IsTty,
	})
	if err != nil {
		return errors.Wrap(err, "error container exec attach")
	}
	defer execAttachResponse.Close()

	handleStreams(opts, &execAttachResponse, opts.OnStreamCloseCallback, opts.OnStreamErrorCallback)

	// fixes echoes and handle SIGTERM interrupt properly
	if terminal, err := util.NewRawTerminal(opts.InStream); err == nil {
		defer terminal.Restore()
	}

	opts.OnContainerAttachCallback()

	// waits for interrupt signals
	statusCh, errCh := client.docker.ContainerWait(client.ctx, opts.ContainerId, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return errors.Wrap(err, "error container wait")
		}
	case <-statusCh:
	}
	return nil
}

func handleStreams(
	opts *ContainerAttachOpts,
	execAttachResponse *types.HijackedResponse,
	onStreamCloseCallback func(),
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

		once.Do(onStreamCloseCallback)
	}()
	go func() {
		if _, err := io.Copy(execAttachResponse.Conn, opts.InStream); err != nil {
			onStreamErrorCallback(errors.Wrap(err, "error copy stdin local->docker"))
		}

		once.Do(onStreamCloseCallback)
	}()
}

func (client *DockerClient) ContainerRestart(containerId string) error {
	if err := client.docker.ContainerRestart(client.ctx, containerId, container.StopOptions{}); err != nil {
		return errors.Wrap(err, "error docker restart")
	}
	return nil
}

func (client *DockerClient) ContainerRemove(containerId string) error {
	if err := client.docker.ContainerRemove(client.ctx, containerId, types.ContainerRemoveOptions{Force: true}); err != nil {
		return errors.Wrap(err, "error docker remove")
	}
	return nil
}

func (client *DockerClient) ContainerList(namePrefix string) ([]ContainerInfo, error) {

	containers, err := client.docker.ContainerList(client.ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key: "name", Value: namePrefix,
		}),
	})
	if err != nil {
		return nil, errors.Wrap(err, "error container list")
	}

	var result []ContainerInfo
	for _, c := range containers {
		// see types.ContainerState
		if c.State == "running" {
			// name starts with slash
			containerName := strings.TrimPrefix(c.Names[0], "/")
			result = append(result, ContainerInfo{ContainerId: c.ID, ContainerName: containerName})
		}
	}

	return result, nil
}
