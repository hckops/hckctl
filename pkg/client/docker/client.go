package docker

import (
	"context"
	"io"
	"strings"
	"sync"

	"github.com/docker/docker/api/types"
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

	// TODO labels
	// TODO env var

	newContainer, err := client.docker.ContainerCreate(
		client.ctx,
		opts.ContainerConfig,
		opts.HostConfig,
		opts.NetworkingConfig,
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

func defaultShell(command string) string {
	if shellCmd := strings.TrimSpace(command); shellCmd != "" {
		return shellCmd
	} else {
		return "/bin/bash"
	}
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

	// fixes echoes and handle SIGTERM interrupt properly
	terminal, err := util.NewRawTerminal(opts.InStream)
	if err != nil {
		return errors.Wrap(err, "error container exec terminal")
	}

	doneChan := make(chan struct{}, 1)
	onStreamCloseCallback := func() {
		terminal.Restore()
		opts.OnStreamCloseCallback()
		close(doneChan)
	}

	handleStreams(opts, &execAttachResponse, onStreamCloseCallback, opts.OnStreamErrorCallback)

	opts.OnContainerAttachCallback()

	// waits for interrupt signals, alternative ContainerExecKill https://github.com/moby/moby/pull/41548
	select {
	case <-client.ctx.Done():
		return client.ctx.Err()
	case <-doneChan:
		return nil
	}
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

func (client *DockerClient) ContainerLogs(opts *ContainerLogsOpts) error {

	outStream, err := client.docker.ContainerLogs(client.ctx, opts.ContainerId, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Follow:     true,
		Details:    true,
	})
	if err != nil {
		return errors.Wrap(err, "error container logs")
	}

	doneChan := make(chan struct{}, 1)
	onStreamCloseCallback := func() {
		// sometimes event is lost if too fast
		opts.OnStreamCloseCallback()
		close(doneChan)
	}

	var once sync.Once
	go func() {
		if _, err := io.Copy(opts.OutStream, outStream); err != nil {
			opts.OnStreamErrorCallback(errors.Wrap(err, "error copy stdout and stderr docker->local"))
		}
		once.Do(onStreamCloseCallback)
	}()

	select {
	case <-client.ctx.Done():
		return client.ctx.Err()
	case <-doneChan:
		return nil
	}
}

func (client *DockerClient) NetworkUpsert(networkName string) (string, error) {

	networks, err := client.docker.NetworkList(client.ctx, types.NetworkListOptions{})
	if err != nil {
		return "", errors.Wrap(err, "error docker network list")
	}
	for _, network := range networks {
		if network.Name == networkName {
			return network.ID, nil
		}
	}

	if newNetwork, err := client.docker.NetworkCreate(client.ctx, networkName, types.NetworkCreate{CheckDuplicate: true}); err != nil {
		return "", errors.Wrap(err, "error docker network create")
	} else {
		return newNetwork.ID, nil
	}
}
