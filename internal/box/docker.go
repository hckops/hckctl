package box

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/rs/zerolog/log"

	"github.com/hckops/hckctl/internal/model"
	"github.com/hckops/hckctl/internal/terminal"
)

type DockerBox struct {
	ctx          context.Context
	dockerClient *client.Client
	loader       *terminal.Loader
	boxTemplate  *model.BoxV1
}

func NewDockerBox(box *model.BoxV1) *DockerBox {

	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal().Err(err).Msg("error docker client")
	}
	defer docker.Close()

	return &DockerBox{
		ctx:          context.Background(),
		dockerClient: docker,
		loader:       terminal.NewLoader(),
		boxTemplate:  box,
	}
}

// TODO add flags detached and tunnel-only
func (b *DockerBox) InitBox() {
	log.Debug().Msgf("init docker box: \n%v\n", b.boxTemplate.Pretty())
	b.loader.Start(fmt.Sprintf("loading %s", b.boxTemplate.Name))

	// TODO compare latest local and remote hash i.e. midnight schedule
	reader, err := b.dockerClient.ImagePull(b.ctx, b.boxTemplate.ImageName(), types.ImagePullOptions{})
	if err != nil {
		log.Fatal().Err(err).Msg("error image pull")
	}
	defer reader.Close()

	b.loader.Refresh(fmt.Sprintf("pulling %s", b.boxTemplate.ImageName()))
	// suppress default output
	io.Copy(ioutil.Discard, reader)

	containerName := b.boxTemplate.GenerateName()

	b.loader.Refresh(fmt.Sprintf("creating %s", containerName))

	// TODO if port is busy start on port+1? or prompt to attach to existing?
	ports := buildDockerPorts(b.boxTemplate.NetworkPorts())

	newContainer, err := b.dockerClient.ContainerCreate(
		b.ctx,
		buildContainerConfig(
			b.boxTemplate.ImageName(),
			containerName,
			ports,
		), // containerConfig
		buildHostConfig(ports), // hostConfig
		nil,                    // networkingConfig
		nil,                    // platform
		containerName)
	if err != nil {
		log.Fatal().Err(err).Msg("error container create")
	}

	containerId := newContainer.ID

	log.Debug().Msgf("open new box: image=%s, containerName=%s, containerId=%s", b.boxTemplate.ImageName(), containerName, containerId)

	// TODO tty false for tunnel only
	b.openBox(containerId, true)
}

func buildDockerPorts(ports []model.PortV1) []nat.Port {

	dockerPorts := make([]nat.Port, 0)
	for _, port := range ports {

		p, err := nat.NewPort("tcp", port.Local)
		if err != nil {
			log.Fatal().Err(err).Msg("error docker port")
		}
		dockerPorts = append(dockerPorts, p)
	}
	return dockerPorts
}

func buildContainerConfig(imageName string, containerName string, ports []nat.Port) *container.Config {

	exposedPorts := make(nat.PortSet)
	for _, port := range ports {
		exposedPorts[port] = struct{}{}
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

func buildHostConfig(ports []nat.Port) *container.HostConfig {

	portBindings := make(nat.PortMap)
	for _, port := range ports {
		portBindings[port] = []nat.PortBinding{{
			HostIP:   "0.0.0.0",
			HostPort: port.Port(),
		}}
	}

	return &container.HostConfig{
		PortBindings: portBindings,
	}
}

func (b *DockerBox) openBox(containerId string, tty bool) {

	if err := b.dockerClient.ContainerStart(b.ctx, containerId, types.ContainerStartOptions{}); err != nil {
		log.Fatal().Err(err).Msg("error container start")
	}

	// TODO always bash
	execCreateResponse, err := b.dockerClient.ContainerExecCreate(b.ctx, containerId, types.ExecConfig{
		AttachStdout: true,
		AttachStdin:  true,
		AttachStderr: true,
		Detach:       false,
		Tty:          tty,
		Cmd:          []string{"/bin/bash"},
	})
	if err != nil {
		log.Fatal().Err(err).Msg("error container exec create")
	}

	execAttachResponse, err := b.dockerClient.ContainerExecAttach(b.ctx, execCreateResponse.ID, types.ExecStartCheck{
		Tty: tty,
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

	handleStreams(&execAttachResponse, tty, removeContainerCallback)

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

func handleStreams(execAttachResponse *types.HijackedResponse, tty bool, onCloseCallback func()) {
	var once sync.Once
	go func() {

		if tty {
			_, err := io.Copy(os.Stdout, execAttachResponse.Reader)
			if err != nil {
				log.Fatal().Err(err).Msg("error copy stdout docker->local")
			}
		} else {
			_, err := stdcopy.StdCopy(os.Stdout, os.Stderr, execAttachResponse.Reader)
			if err != nil {
				log.Fatal().Err(err).Msg("error copy stdout and stderr docker->local")
			}
		}

		once.Do(onCloseCallback)
	}()
	go func() {
		_, err := io.Copy(execAttachResponse.Conn, os.Stdin)
		if err != nil {
			log.Fatal().Err(err).Msg("error copy stdin local->docker")
		}

		once.Do(onCloseCallback)
	}()
}
