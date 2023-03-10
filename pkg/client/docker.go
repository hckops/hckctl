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
	"github.com/moby/term"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log" // TODO remove

	"github.com/hckops/hckctl/pkg/model"
	"github.com/hckops/hckctl/pkg/util"
)

type DockerBox struct {
	ctx      context.Context
	docker   *client.Client
	Template *model.BoxV1
}

// TODO don't forget to invoke "dockerClient.Close()"
func NewDockerBox(template *model.BoxV1) (*DockerBox, error) {

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

// TODO remove onSetupCallback, invoke before/after
func (box *DockerBox) Setup(onSetupCallback func()) error {

	// TODO delete dangling images
	reader, err := box.docker.ImagePull(box.ctx, box.Template.ImageName(), types.ImagePullOptions{})
	if err != nil {
		return errors.Wrap(err, "error image pull")
	}
	defer reader.Close()

	onSetupCallback()

	// suppress default output
	if _, err := io.Copy(ioutil.Discard, reader); err != nil {
		return errors.Wrap(err, "error image pull output message")
	}

	return nil
}

func (box *DockerBox) Create(containerName string, onPortBindCallback func(port model.PortV1)) (string, error) {

	containerConfig, err := buildContainerConfig(
		box.Template.ImageName(),
		containerName,
		box.Template.NetworkPorts(),
	)
	if err != nil {
		return "", err
	}

	hostConfig, err := buildHostConfig(box.Template.NetworkPorts(), onPortBindCallback)
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

func buildContainerConfig(imageName string, containerName string, ports []model.PortV1) (*container.Config, error) {

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

func buildHostConfig(ports []model.PortV1, onPortBindCallback func(port model.PortV1)) (*container.HostConfig, error) {

	portBindings := make(nat.PortMap)
	for _, port := range ports {
		remotePort, err := nat.NewPort("tcp", port.Remote)
		if err != nil {
			return nil, errors.Wrap(err, "error docker port: hostConfig")
		}

		localPort := util.GetLocalPort(port.Local)

		// actual port bindedmodel
		onPortBindCallback(model.PortV1{
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

// TODO set bash in the config
func (box *DockerBox) Exec(containerId string, streams *model.BoxStreams, onExecCallback func()) error {

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

	// TODO >>> LOG
	removeContainerCallback := func() {
		log.Debug().Msgf("removing container: %s", containerId)
		if err := box.docker.ContainerRemove(box.ctx, containerId, types.ContainerRemoveOptions{Force: true}); err != nil {
			log.Fatal().Err(err).Msg("error docker remove")
		}
	}

	handleStreams(&execAttachResponse, streams, removeContainerCallback)

	// fixes echoes and handle SIGTERM interrupt properly
	if fd, isTerminal := term.GetFdInfo(streams.Stdin); isTerminal {
		previousState, err := term.SetRawTerminal(fd)
		if err != nil {
			return errors.Wrap(err, "error raw terminal")
		}
		defer term.RestoreTerminal(fd, previousState)
	}

	onExecCallback()

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

// TODO >>> LOG
func handleStreams(execAttachResponse *types.HijackedResponse, streams *model.BoxStreams, onCloseCallback func()) {
	var once sync.Once
	go func() {

		if streams.IsTty {
			_, err := io.Copy(streams.Stdout, execAttachResponse.Reader)
			if err != nil {
				log.Fatal().Err(err).Msg("error copy stdout docker->local")
			}
		} else {
			_, err := stdcopy.StdCopy(streams.Stdout, streams.Stderr, execAttachResponse.Reader)
			if err != nil {
				log.Fatal().Err(err).Msg("error copy stdout and stderr docker->local")
			}
		}

		once.Do(onCloseCallback)
	}()
	go func() {
		_, err := io.Copy(execAttachResponse.Conn, streams.Stdin)
		if err != nil {
			log.Fatal().Err(err).Msg("error copy stdin local->docker")
		}

		once.Do(onCloseCallback)
	}()
}
