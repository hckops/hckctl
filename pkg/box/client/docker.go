package client

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/template/model"
	"github.com/hckops/hckctl/pkg/util"
)

type DockerBox struct {
	ctx      context.Context
	docker   *client.Client
	template *model.BoxV1

	OutChan chan string
	ErrChan chan string
}

// TODO don't forget to invoke "defer dockerClient.Close()"

func NewDockerBox(template *model.BoxV1) (*DockerBox, error) {

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, errors.Wrap(err, "error docker client")
	}

	return &DockerBox{
		ctx:      context.Background(),
		docker:   dockerClient,
		template: template,
	}, nil
}

func (box *DockerBox) Close() error {
	return box.docker.Close()
}

func (box *DockerBox) Setup() error {

	// TODO delete dangling images

	box.OutChan <- "setup"

	reader, err := box.docker.ImagePull(box.ctx, box.template.ImageName(), types.ImagePullOptions{})
	if err != nil {
		return errors.Wrap(err, "error image pull")
	}
	defer reader.Close()

	// TODO
	box.OutChan <- "setup-image-pull"

	// suppress default output
	if _, err := io.Copy(io.Discard, reader); err != nil {
		return errors.Wrap(err, "error image pull output message")
	}

	box.OutChan <- "setup-copy"

	return nil
}

func (box *DockerBox) Create(containerName string) (string, error) {

	containerConfig, err := buildContainerConfig(
		box.template.ImageName(),
		containerName,
		box.template.NetworkPorts(),
	)
	if err != nil {
		return "", err
	}

	hostConfig, err := buildHostConfig(box.template.NetworkPorts(), box.OutChan)
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

func buildContainerConfig(imageName string, containerName string, ports []model.BoxPort) (*container.Config, error) {

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

func buildHostConfig(ports []model.BoxPort, outChan chan string) (*container.HostConfig, error) {

	portBindings := make(nat.PortMap)
	for _, port := range ports {

		localPort, err := util.FindOpenPort(port.Local)
		if err != nil {
			return nil, errors.Wrap(err, "error docker local port: hostConfig")
		}

		remotePort, err := nat.NewPort("tcp", port.Remote)
		if err != nil {
			return nil, errors.Wrap(err, "error docker remote port: hostConfig")
		}

		// TODO
		// actual binded port
		_ = model.BoxPort{
			Alias:  port.Alias,
			Local:  localPort,
			Remote: port.Remote,
		}
		// TODO
		outChan <- "onPortBindCallback"

		portBindings[remotePort] = []nat.PortBinding{{
			HostIP:   "0.0.0.0",
			HostPort: localPort,
		}}
	}

	return &container.HostConfig{
		PortBindings: portBindings,
	}, nil
}
