package box

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/hckops/hckctl/pkg/util"
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/template/model"
)

type DockerClient struct {
	ctx       context.Context
	dockerApi *client.Client
	template  *model.BoxV1
	eventChan chan BoxEvent
	wg        sync.WaitGroup
}

func NewDockerClient(template *model.BoxV1) (*DockerClient, error) {

	dockerApi, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, errors.Wrap(err, "error docker client")
	}

	return &DockerClient{
		ctx:       context.Background(),
		dockerApi: dockerApi,
		template:  template,
		eventChan: make(chan BoxEvent),
	}, nil
}

func (c *DockerClient) Events() <-chan BoxEvent {
	return c.eventChan
}

func (c *DockerClient) sendDebugEvent(source, message string) {
	c.wg.Add(1)
	go func() {
		c.eventChan <- BoxEvent{
			Kind:    DebugEvent,
			Source:  source,
			Message: message,
		}
		c.wg.Done()
	}()
}

func (c *DockerClient) sendSuccessEvent(source string) {
	c.wg.Add(1)
	go func() {
		c.eventChan <- BoxEvent{
			Kind:    SuccessEvent,
			Source:  source,
			Message: "",
		}
		c.wg.Done()
	}()
}

func (c *DockerClient) close() error {
	close(c.eventChan)
	return c.dockerApi.Close()
}

func (c *DockerClient) Create() (string, error) {
	defer c.close()
	if err := c.setup(); err != nil {
		// TODO
		return "", err
	}

	boxName := c.template.GenerateName()
	boxId, err := c.createContainer(boxName)
	if err != nil {
		// TODO
		return "", err
	}

	c.wg.Wait()

	return boxId, nil
}

func (c *DockerClient) setup() error {

	// TODO delete dangling images

	c.sendDebugEvent("setup", "step-1")

	reader, err := c.dockerApi.ImagePull(c.ctx, c.template.ImageName(), types.ImagePullOptions{})
	if err != nil {
		return errors.Wrap(err, "error image pull")
	}
	defer reader.Close()

	// TODO
	c.sendDebugEvent("setup", "step-2")

	// suppress default output
	if _, err := io.Copy(io.Discard, reader); err != nil {
		return errors.Wrap(err, "error image pull output message")
	}

	c.sendDebugEvent("setup", "step-3")

	return nil
}

func (c *DockerClient) createContainer(containerName string) (string, error) {

	c.sendDebugEvent("create", "step-1")

	containerConfig, err := buildContainerConfig(
		c.template.ImageName(),
		containerName,
		c.template.NetworkPorts(),
	)
	if err != nil {
		return "", err
	}

	onPortBindCallback := func(port model.BoxPort) {
		c.sendDebugEvent("create-port", fmt.Sprintf("port %v", port))
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

	c.sendSuccessEvent("create")

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
