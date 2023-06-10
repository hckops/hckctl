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

type DockerBoxOld struct {
	ctx        context.Context
	docker     *client.Client
	template   *model.BoxV1
	statusChan chan string
	errorChan  chan error
}

func NewDockerBoxOld(template *model.BoxV1, statusChan chan string, errorChan chan error) (*DockerBoxOld, error) {

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, errors.Wrap(err, "error docker client")
	}

	return &DockerBoxOld{
		ctx:        context.Background(),
		docker:     dockerClient,
		template:   template,
		statusChan: statusChan,
		errorChan:  errorChan,
	}, nil
}

// TODO don't forget to invoke "defer dockerClient.Close()"

func (box *DockerBoxOld) updateStatus(value string) {
	go func() {
		box.statusChan <- value
	}()
}

func (box *DockerBoxOld) Close() error {
	close(box.statusChan)
	close(box.errorChan)
	return box.docker.Close()
}

func (box *DockerBoxOld) Setup() error {

	// TODO delete dangling images

	box.updateStatus("setup")

	reader, err := box.docker.ImagePull(box.ctx, box.template.ImageName(), types.ImagePullOptions{})
	if err != nil {
		return errors.Wrap(err, "error image pull")
	}
	defer reader.Close()

	// TODO
	box.updateStatus("setup-image-pull")

	// suppress default output
	if _, err := io.Copy(io.Discard, reader); err != nil {
		return errors.Wrap(err, "error image pull output message")
	}

	box.updateStatus("setup-copy")

	return nil
}

func (box *DockerBoxOld) Create(containerName string) (string, error) {

	box.updateStatus("create")

	containerConfig, err := buildContainerConfig(
		box.template.ImageName(),
		containerName,
		box.template.NetworkPorts(),
	)
	if err != nil {
		return "", err
	}

	hostConfig, err := buildHostConfig(box.template.NetworkPorts(), box.statusChan)
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

	go func() {
		box.errorChan <- nil
	}()
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
		go func() {
			outChan <- "onPortBindCallback"
		}()

		portBindings[remotePort] = []nat.PortBinding{{
			HostIP:   "0.0.0.0",
			HostPort: localPort,
		}}
	}

	return &container.HostConfig{
		PortBindings: portBindings,
	}, nil
}
