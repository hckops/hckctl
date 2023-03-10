package docker

import (
	"io"
	"io/ioutil"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/rs/zerolog/log"

	"github.com/hckops/hckctl/internal/terminal" // TODO
	"github.com/hckops/hckctl/pkg/model"
	"github.com/hckops/hckctl/pkg/util"
)

type DockerBox struct {
	DockerClient *client.Client
	Context      *model.BoxContext
}

func (box *DockerBox) Setup(onSetupCallback func()) {

	// TODO delete dangling images
	reader, err := box.DockerClient.ImagePull(box.Context.Ctx, box.Context.Template.ImageName(), types.ImagePullOptions{})
	if err != nil {
		log.Fatal().Err(err).Msg("error image pull")
	}
	defer reader.Close()

	onSetupCallback()

	// suppress default output
	io.Copy(ioutil.Discard, reader)
}

func (box *DockerBox) Create(containerName string) string {

	newContainer, err := box.DockerClient.ContainerCreate(
		box.Context.Ctx,
		buildContainerConfig(
			box.Context.Template.ImageName(),
			containerName,
			box.Context.Template.NetworkPorts(),
		), // containerConfig
		buildHostConfig(box.Context.Template.NetworkPorts()), // hostConfig
		nil, // networkingConfig
		nil, // platform
		containerName)
	if err != nil {
		log.Fatal().Err(err).Msg("error container create")
	}

	return newContainer.ID
}

func buildContainerConfig(imageName string, containerName string, ports []model.PortV1) *container.Config {

	exposedPorts := make(nat.PortSet)
	for _, port := range ports {
		p, err := nat.NewPort("tcp", port.Remote)
		if err != nil {
			log.Fatal().Err(err).Msg("error docker port: containerConfig")
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
	}
}

func buildHostConfig(ports []model.PortV1) *container.HostConfig {

	portBindings := make(nat.PortMap)
	for _, port := range ports {
		remotePort, err := nat.NewPort("tcp", port.Remote)
		if err != nil {
			log.Fatal().Err(err).Msg("error docker port: hostConfig")
		}

		localPort := util.GetLocalPort(port.Local)
		log.Info().Msgf("[%s] exposing %s (local) -> %s (container)", port.Alias, localPort, port.Remote)

		portBindings[remotePort] = []nat.PortBinding{{
			HostIP:   "0.0.0.0",
			HostPort: localPort,
		}}
	}

	return &container.HostConfig{
		PortBindings: portBindings,
	}
}

func (box *DockerBox) Exec(containerId string, onExecCallback func()) {

	if err := box.DockerClient.ContainerStart(box.Context.Ctx, containerId, types.ContainerStartOptions{}); err != nil {
		log.Fatal().Err(err).Msg("error container start")
	}

	// TODO always bash
	execCreateResponse, err := box.DockerClient.ContainerExecCreate(box.Context.Ctx, containerId, types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Detach:       false,
		Tty:          box.Context.Streams.IsTty,
		Cmd:          []string{"/bin/bash"},
	})
	if err != nil {
		log.Fatal().Err(err).Msg("error container exec create")
	}

	execAttachResponse, err := box.DockerClient.ContainerExecAttach(box.Context.Ctx, execCreateResponse.ID, types.ExecStartCheck{
		Tty: box.Context.Streams.IsTty,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("error container exec attach")
	}
	defer execAttachResponse.Close()

	removeContainerCallback := func() {
		if err := box.DockerClient.ContainerRemove(box.Context.Ctx, containerId, types.ContainerRemoveOptions{Force: true}); err != nil {
			log.Fatal().Err(err).Msg("error docker remove")
		}
	}

	handleStreams(&execAttachResponse, box.Context.Streams, removeContainerCallback)

	// fixes echoes and handle SIGTERM interrupt properly
	rawTerminal := terminal.NewRawTerminal()
	if rawTerminal == nil {
		log.Fatal().Msg("error raw terminal")
	}
	defer rawTerminal.Restore()

	onExecCallback()

	// waits for interrupt signals
	statusCh, errCh := box.DockerClient.ContainerWait(box.Context.Ctx, containerId, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			log.Fatal().Err(err).Msg("error container wait")
		}
	case <-statusCh:
	}
}

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
