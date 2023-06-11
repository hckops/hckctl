package box

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/client"
	"github.com/hckops/hckctl/pkg/client/docker"
	"github.com/hckops/hckctl/pkg/template/model"
	"github.com/hckops/hckctl/pkg/util"
)

type DockerBox struct {
	client *docker.DockerClient
	opts   *boxOpts
}

func NewDockerBox(opts *boxOpts) (*DockerBox, error) {

	dockerClient, err := docker.NewDockerClient(opts.eventBus)
	if err != nil {
		return nil, errors.Wrap(err, "error docker box")
	}

	return &DockerBox{
		client: dockerClient,
		opts:   opts,
	}, nil
}

func (b *DockerBox) Events() *client.EventBus {
	return b.opts.eventBus
}

func (b *DockerBox) Create() (*BoxInfo, error) {
	defer b.client.Close()
	return b.createBox()
}

func (b *DockerBox) createBox() (*BoxInfo, error) {

	imageName := b.opts.template.ImageName()
	if err := b.client.Setup(imageName); err != nil {
		return nil, err
	}

	// boxName
	containerName := b.opts.template.GenerateName()
	containerConfig, err := buildContainerConfig(
		b.opts.template.ImageName(),
		containerName,
		b.opts.template.NetworkPorts(),
	)
	if err != nil {
		return nil, err
	}

	onPortBindCallback := func(port model.BoxPort) {
		// TODO generic or docker?
		//c.eventBus.PublishConsoleEvent("createContainer",
		//	"[%s][%s]   \texpose (container) %s -> (local) http://localhost:%s",
		//	containerName, port.Alias, port.Remote, port.Local)
	}

	hostConfig, err := buildHostConfig(b.opts.template.NetworkPorts(), onPortBindCallback)
	if err != nil {
		return nil, err
	}

	createContainerOpts := &docker.CreateContainerOpts{
		ContainerName:   containerName,
		ContainerConfig: containerConfig,
		HostConfig:      hostConfig,
	}

	// boxId
	containerId, err := b.client.CreateContainer(createContainerOpts)
	if err != nil {
		return nil, err
	}
	// TODO generic or docker?
	//c.eventBus.PublishDebugEvent("create", "create box: templateName=%s boxName=%s boxId=%s", c.template.Name, boxName, boxId)

	return &BoxInfo{Id: containerId, Name: containerName}, nil
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

		// actual bound port
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

func (b *DockerBox) Exec(info *BoxInfo) error {

	execContainerOpts := &docker.ExecContainerOpts{
		ContainerId: info.Id,
		Shell:       b.opts.template.Shell,
		InStream:    b.opts.streams.in,
		OutStream:   b.opts.streams.out,
		ErrStream:   b.opts.streams.err,
		IsTty:       b.opts.streams.isTty,
	}

	return b.client.ExecContainer(execContainerOpts)
}

func (b *DockerBox) Copy(info *BoxInfo, from string, to string) error {
	return nil
}

func (b *DockerBox) List() ([]string, error) {
	return nil, nil
}

func (b *DockerBox) Open() error {
	defer b.client.Close()

	info, err := b.createBox()
	if err != nil {
		return err
	}

	return b.Exec(info)
}

func (b *DockerBox) Tunnel() error {
	return nil
}

func (b *DockerBox) Delete(boxId string) error {
	return nil
}
