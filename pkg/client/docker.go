package client

import (
	"context"
	"io"
	"io/ioutil"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/model"
	"github.com/hckops/hckctl/pkg/schema"
	"github.com/hckops/hckctl/pkg/util"
)

type DockerBox struct {
	ctx      context.Context
	docker   *client.Client
	Template *schema.BoxV1

	// TODO look for better design than callbacks
	OnSetupCallback       func()
	OnCreateCallback      func(port schema.PortV1)
	OnExecCallback        func()
	OnCloseCallback       func()
	OnCloseErrorCallback  func(error, string)
	OnStreamErrorCallback func(error, string) // TODO replace with error channel
}

// don't forget to invoke "defer dockerClient.Close()"
func NewDockerBox(template *schema.BoxV1) (*DockerBox, error) {

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, errors.Wrap(err, "error docker client")
	}

	return &DockerBox{
		ctx:      context.Background(),
		docker:   dockerClient,
		Template: template,
	}, nil
}

func (box *DockerBox) Close() error {
	return box.docker.Close()
}

func (box *DockerBox) Setup() error {

	// TODO delete dangling images
	reader, err := box.docker.ImagePull(box.ctx, box.Template.ImageName(), types.ImagePullOptions{})
	if err != nil {
		return errors.Wrap(err, "error image pull")
	}
	defer reader.Close()

	box.OnSetupCallback()

	// suppress default output
	if _, err := io.Copy(ioutil.Discard, reader); err != nil {
		return errors.Wrap(err, "error image pull output message")
	}

	return nil
}

func (box *DockerBox) Create(containerName string) (string, error) {

	containerConfig, err := buildContainerConfig(
		box.Template.ImageName(),
		containerName,
		box.Template.NetworkPorts(),
	)
	if err != nil {
		return "", err
	}

	hostConfig, err := buildHostConfig(box.Template.NetworkPorts(), box.OnCreateCallback)
	if err != nil {
		return "", err
	}

	newContainer, err := box.docker.ContainerCreate(
		box.ctx,
		containerConfig,
		hostConfig,
		nil, // networkingConfig
		nil, // platform
		containerName)
	if err != nil {
		return "", errors.Wrap(err, "error container create")
	}

	return newContainer.ID, nil
}

func buildContainerConfig(imageName string, containerName string, ports []schema.PortV1) (*container.Config, error) {

	exposedPorts := make(nat.PortSet)
	for _, port := range ports {
		p, err := nat.NewPort("tcp", port.Remote)
		if err != nil {
			return nil, errors.Wrap(err, "error docker port: containerConfig")
		}
		exposedPorts[p] = struct{}{}
	}

	return &container.Config{
		Hostname:     containerName,
		Image:        imageName,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		OpenStdin:    true,
		StdinOnce:    true,
		Tty:          true,
		ExposedPorts: exposedPorts,
	}, nil
}

func buildHostConfig(ports []schema.PortV1, onPortBindCallback func(port schema.PortV1)) (*container.HostConfig, error) {

	portBindings := make(nat.PortMap)
	for _, port := range ports {

		localPort, err := util.GetLocalPort(port.Local)
		if err != nil {
			return nil, errors.Wrap(err, "error docker local port: hostConfig")
		}

		remotePort, err := nat.NewPort("tcp", port.Remote)
		if err != nil {
			return nil, errors.Wrap(err, "error docker remote port: hostConfig")
		}

		// actual binded port
		onPortBindCallback(schema.PortV1{
			Alias:  port.Alias,
			Local:  localPort,
			Remote: port.Remote,
		})

		portBindings[remotePort] = []nat.PortBinding{{
			HostIP:   "0.0.0.0",
			HostPort: localPort,
		}}
	}

	return &container.HostConfig{
		PortBindings: portBindings,
	}, nil
}

// TODO command from template
func (box *DockerBox) Exec(containerId string, streams *model.BoxStreams) error {

	if err := box.docker.ContainerStart(box.ctx, containerId, types.ContainerStartOptions{}); err != nil {
		return errors.Wrap(err, "error container start")
	}

	execCreateResponse, err := box.docker.ContainerExecCreate(box.ctx, containerId, types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Detach:       false,
		Tty:          streams.IsTty,
		Cmd:          []string{"/bin/bash"},
	})
	if err != nil {
		return errors.Wrap(err, "error container exec create")
	}

	execAttachResponse, err := box.docker.ContainerExecAttach(box.ctx, execCreateResponse.ID, types.ExecStartCheck{
		Tty: streams.IsTty,
	})
	if err != nil {
		return errors.Wrap(err, "error container exec attach")
	}
	defer execAttachResponse.Close()

	removeContainerCallback := func() {
		box.OnCloseCallback()

		if err := box.docker.ContainerRemove(box.ctx, containerId, types.ContainerRemoveOptions{Force: true}); err != nil {
			box.OnCloseErrorCallback(err, "error docker remove")
		}
	}

	handleStreams(&execAttachResponse, streams, removeContainerCallback, box.OnStreamErrorCallback)

	// fixes echoes and handle SIGTERM interrupt properly
	if terminal, err := util.NewRawTerminal(streams.Stdin); err == nil {
		defer terminal.Restore()
	}

	box.OnExecCallback()

	// waits for interrupt signals
	statusCh, errCh := box.docker.ContainerWait(box.ctx, containerId, container.WaitConditionNotRunning)
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
	execAttachResponse *types.HijackedResponse,
	streams *model.BoxStreams,
	onCloseCallback func(),
	onStreamErrorCallback func(error, string),
) {

	var once sync.Once
	go func() {

		if streams.IsTty {
			if _, err := io.Copy(streams.Stdout, execAttachResponse.Reader); err != nil {
				onStreamErrorCallback(err, "error copy stdout docker->local")
			}
		} else {
			if _, err := stdcopy.StdCopy(streams.Stdout, streams.Stderr, execAttachResponse.Reader); err != nil {
				onStreamErrorCallback(err, "error copy stdout and stderr docker->local")
			}
		}

		once.Do(onCloseCallback)
	}()
	go func() {
		if _, err := io.Copy(execAttachResponse.Conn, streams.Stdin); err != nil {
			onStreamErrorCallback(err, "error copy stdin local->docker")
		}

		once.Do(onCloseCallback)
	}()
}
