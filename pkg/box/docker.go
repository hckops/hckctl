package box

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

type DockerClient struct {
	ctx       context.Context
	dockerApi *client.Client
	template  *model.BoxV1
	eventBus  *EventBus
}

func NewDockerClient(template *model.BoxV1, eventBus *EventBus) (*DockerClient, error) {

	dockerApi, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, errors.Wrap(err, "error docker client")
	}

	return &DockerClient{
		ctx:       context.Background(),
		dockerApi: dockerApi,
		template:  template,
		eventBus:  eventBus,
	}, nil
}

func (c *DockerClient) Events() *EventBus {
	return c.eventBus
}

func (c *DockerClient) Create() (string, error) {
	defer c.close()

	// TODO delete dangling images

	if err := c.setup(); err != nil {
		return "", err
	}

	boxName := c.template.GenerateName()
	boxId, err := c.createContainer(boxName)
	if err != nil {
		return "", err
	}
	c.eventBus.PublishInfoEvent("create", "create box: templateName=%s boxName=%s boxId=%s", c.template.Name, boxName, boxId)

	return boxName, nil
}

func (c *DockerClient) close() error {
	c.eventBus.PublishDebugEvent("close", "closing event bus and docker client")
	c.eventBus.Close()
	return c.dockerApi.Close()
}

func (c *DockerClient) setup() error {
	imageName := c.template.ImageName()

	reader, err := c.dockerApi.ImagePull(c.ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return errors.Wrap(err, "error image pull")
	}
	defer reader.Close()

	c.eventBus.PublishPriorityEvent("setup", "pulling %s", imageName)

	// suppress default output
	if _, err := io.Copy(io.Discard, reader); err != nil {
		return errors.Wrap(err, "error image pull output message")
	}

	return nil
}

func (c *DockerClient) createContainer(containerName string) (string, error) {

	containerConfig, err := buildContainerConfig(
		c.template.ImageName(),
		containerName,
		c.template.NetworkPorts(),
	)
	if err != nil {
		return "", err
	}

	onPortBindCallback := func(port model.BoxPort) {
		c.eventBus.PublishInfoEvent("createContainer",
			"[%s][%s] exposing %s (local) -> %s (container)",
			containerName, port.Alias, port.Local, port.Remote)
	}

	hostConfig, err := buildHostConfig(c.template.NetworkPorts(), onPortBindCallback)
	if err != nil {
		return "", err
	}

	newContainer, err := c.dockerApi.ContainerCreate(
		c.ctx,
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

func buildHostConfig(ports []model.BoxPort, onPortBindCallback func(port model.BoxPort)) (*container.HostConfig, error) {

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

		// actual binded port
		onPortBindCallback(model.BoxPort{
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
