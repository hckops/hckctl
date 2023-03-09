package box

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/rs/zerolog/log"

	"github.com/hckops/hckctl/internal/common"
	"github.com/hckops/hckctl/internal/model"
	"github.com/hckops/hckctl/internal/terminal"
)

type DockerBox struct {
	ctx          context.Context
	loader       *terminal.Loader
	template     *model.BoxV1
	dockerClient *client.Client
}

func NewDockerBox(template *model.BoxV1) *DockerBox {

	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal().Err(err).Msg("error docker client")
	}
	defer docker.Close()

	return &DockerBox{
		ctx:          context.Background(),
		loader:       terminal.NewLoader(),
		template:     template,
		dockerClient: docker,
	}
}

func (b *DockerBox) OpenBox(streams *model.BoxStreams) {
	log.Debug().Msgf("init docker box: \n%v\n", b.template.Pretty())
	b.loader.Start(fmt.Sprintf("loading %s", b.template.Name))

	// TODO delete dangling images
	reader, err := b.dockerClient.ImagePull(b.ctx, b.template.ImageName(), types.ImagePullOptions{})
	if err != nil {
		log.Fatal().Err(err).Msg("error image pull")
	}
	defer reader.Close()

	b.loader.Refresh(fmt.Sprintf("pulling %s", b.template.ImageName()))
	// suppress default output
	io.Copy(ioutil.Discard, reader)

	containerName := b.template.GenerateName()

	b.loader.Refresh(fmt.Sprintf("creating %s", containerName))

	newContainer, err := b.dockerClient.ContainerCreate(
		b.ctx,
		buildContainerConfig(
			b.template.ImageName(),
			containerName,
			b.template.NetworkPorts(),
		), // containerConfig
		buildHostConfig(b.template.NetworkPorts()), // hostConfig
		nil, // networkingConfig
		nil, // platform
		containerName)
	if err != nil {
		log.Fatal().Err(err).Msg("error container create")
	}

	containerId := newContainer.ID

	log.Info().Msgf("open new box: image=%s, containerName=%s, containerId=%s", b.template.ImageName(), containerName, containerId)

	b.execContainer(containerId, streams)
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

		localPort := common.GetLocalPort(port.Local)
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

func (b *DockerBox) execContainer(containerId string, streams *model.BoxStreams) {

	if err := b.dockerClient.ContainerStart(b.ctx, containerId, types.ContainerStartOptions{}); err != nil {
		log.Fatal().Err(err).Msg("error container start")
	}

	// TODO always bash
	execCreateResponse, err := b.dockerClient.ContainerExecCreate(b.ctx, containerId, types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Detach:       false,
		Tty:          streams.IsTty,
		Cmd:          []string{"/bin/bash"},
	})
	if err != nil {
		log.Fatal().Err(err).Msg("error container exec create")
	}

	execAttachResponse, err := b.dockerClient.ContainerExecAttach(b.ctx, execCreateResponse.ID, types.ExecStartCheck{
		Tty: streams.IsTty,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("error container exec attach")
	}
	defer execAttachResponse.Close()

	removeContainerCallback := func() {
		if err := b.dockerClient.ContainerRemove(b.ctx, containerId, types.ContainerRemoveOptions{Force: true}); err != nil {
			log.Fatal().Err(err).Msg("error docker remove")
		}
	}

	handleStreams(&execAttachResponse, streams, removeContainerCallback)

	// fixes echoes and handle SIGTERM interrupt properly
	rawTerminal := terminal.NewRawTerminal()
	if rawTerminal == nil {
		log.Fatal().Msg("error raw terminal")
	}
	defer rawTerminal.Restore()
	b.loader.Stop()

	// waits for interrupt signals
	statusCh, errCh := b.dockerClient.ContainerWait(b.ctx, containerId, container.WaitConditionNotRunning)
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
