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
	opts.eventBus.Publish(newInitBoxEvent())

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

func (b *DockerBox) Create(template *model.BoxV1) (*BoxInfo, error) {
	defer b.client.Close()
	return b.createBox(template)
}

// TODO exclude tty experimental port
func (b *DockerBox) createBox(template *model.BoxV1) (*BoxInfo, error) {

	imageName := template.ImageName()
	setupImageOpts := &docker.SetupImageOpts{
		ImageName: imageName,
		OnPullImageCallback: func() {
			b.opts.eventBus.Publish(newPullImageBoxEvent(imageName))
		},
	}
	if err := b.client.Setup(setupImageOpts); err != nil {
		return nil, err
	}

	// boxName
	containerName := template.GenerateName()
	containerConfig, err := buildContainerConfig(
		template.ImageName(),
		containerName,
		template.NetworkPorts(),
	)
	if err != nil {
		return nil, err
	}

	onPortBindCallback := func(port model.BoxPort) {
		b.opts.eventBus.Publish(newBindPortBoxEvent(containerName, port))
	}

	hostConfig, err := buildHostConfig(template.NetworkPorts(), onPortBindCallback)
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

	b.opts.eventBus.Publish(newGenericBoxEvent("new box created successfully: templateName=%s boxName=%s boxId=%s",
		template.Name, containerName, containerId))

	return &BoxInfo{Id: containerId, Name: containerName, Shell: template.Shell}, nil
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

func (b *DockerBox) Exec(info BoxInfo) error {

	// TODO resolve id by name

	execContainerOpts := &docker.ExecContainerOpts{
		ContainerId: info.Id,
		Shell:       info.Shell,
		InStream:    b.opts.streams.in,
		OutStream:   b.opts.streams.out,
		ErrStream:   b.opts.streams.err,
		IsTty:       b.opts.streams.isTty,
		OnContainerWaitingCallback: func() {
			b.opts.eventBus.Publish(newContainerWaitingBoxEvent())
		},
	}

	return b.client.ExecContainer(execContainerOpts)
}

func (b *DockerBox) Copy(info BoxInfo, from string, to string) error {
	return nil
}

func (b *DockerBox) List() ([]BoxInfo, error) {
	defer b.client.Close()

	containers, err := b.client.ListContainers()
	if err != nil {
		return nil, err
	}
	var result []BoxInfo
	for _, c := range containers {
		result = append(result, BoxInfo{Id: c.ContainerId, Name: c.ContainerName})
	}

	return result, nil
}

func (b *DockerBox) Open(template *model.BoxV1) error {
	defer b.client.Close()

	info, err := b.createBox(template)
	if err != nil {
		return err
	}

	return b.Exec(*info)
}

func (b *DockerBox) Tunnel(BoxInfo) error {
	return errors.New("not supported")
}

func (b *DockerBox) Delete(info BoxInfo) error {
	defer b.client.Close()
	// TODO resolve id by name
	return b.client.RemoveContainer(info.Id)
}
