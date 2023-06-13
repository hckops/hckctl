package docker

import (
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/box"
	"github.com/hckops/hckctl/pkg/client/docker"
	"github.com/hckops/hckctl/pkg/template/model"
	"github.com/hckops/hckctl/pkg/util"
)

const (
	labelRevision = "com.hckops.revision"
)

type DockerBox struct {
	client *docker.DockerClient
	opts   *box.BoxOpts
}

func NewDockerBox(opts *box.BoxOpts) (*DockerBox, error) {
	opts.EventBus.Publish(newClientInitDockerEvent())

	dockerClient, err := docker.NewDockerClient()
	if err != nil {
		return nil, errors.Wrap(err, "error docker box")
	}

	return &DockerBox{
		client: dockerClient,
		opts:   opts,
	}, nil
}

func (b *DockerBox) Events() *box.EventBus {
	return b.opts.EventBus
}

func (b *DockerBox) close() error {
	b.opts.EventBus.Publish(newClientCloseDockerEvent())
	b.opts.EventBus.Close()
	return b.client.Close()
}

func (b *DockerBox) Create(template *model.BoxV1) (*box.BoxInfo, error) {
	defer b.close()
	return b.createBox(template)
}

func (b *DockerBox) createBox(template *model.BoxV1) (*box.BoxInfo, error) {

	imageName := template.ImageName()
	imagePullOpts := &docker.ImagePullOpts{
		ImageName: imageName,
		OnImagePullCallback: func() {
			b.opts.EventBus.Publish(newImagePullDockerLoaderEvent(imageName))
		},
	}
	b.opts.EventBus.Publish(newImagePullDockerEvent(imageName))
	if err := b.client.ImagePull(imagePullOpts); err != nil {
		return nil, err
	}

	// cleanup old images
	imageRemoveOpts := &docker.ImageRemoveOpts{
		OnImageRemoveCallback: func(imageId string) {
			b.opts.EventBus.Publish(newImageRemoveDockerEvent(imageId))
		},
		OnImageRemoveErrorCallback: func(imageId string, err error) {
			b.opts.EventBus.Publish(newImageRemoveErrorDockerEvent(imageId, err))
		},
	}
	if err := b.client.ImageRemoveDangling(imageRemoveOpts); err != nil {
		return nil, err
	}

	// boxName
	containerName := template.GenerateName()
	// skip not supported virtual-* ports
	var networkPorts []model.BoxPort
	for _, networkPort := range template.NetworkPorts() {
		if strings.HasPrefix(networkPort.Alias, model.BoxPrefixVirtualPort) {
			b.opts.EventBus.Publish(newContainerCreateSkipVirtualPortDockerEvent(containerName, networkPort))
		} else {
			networkPorts = append(networkPorts, networkPort)
		}
	}
	containerConfig, err := buildContainerConfig(
		template.ImageName(),
		containerName,
		networkPorts,
	)
	if err != nil {
		return nil, err
	}

	onPortBindCallback := func(port model.BoxPort) {
		b.opts.EventBus.Publish(newContainerCreatePortBindDockerEvent(containerName, port))
		b.opts.EventBus.Publish(newContainerCreatePortBindDockerConsoleEvent(containerName, port))
	}
	hostConfig, err := buildHostConfig(networkPorts, onPortBindCallback)
	if err != nil {
		return nil, err
	}

	containerOpts := &docker.ContainerCreateOpts{
		ContainerName:   containerName,
		ContainerConfig: containerConfig,
		HostConfig:      hostConfig,
	}
	// boxId
	containerId, err := b.client.ContainerCreate(containerOpts)
	if err != nil {
		return nil, err
	}
	b.opts.EventBus.Publish(newContainerCreateDockerEvent(template.Name, containerName, containerId))

	return &box.BoxInfo{Id: containerId, Name: containerName}, nil
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
		Labels:       map[string]string{labelRevision: "TODO-main-or-empty"}, // TODO use correct revision when resolving template by name
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

func (b *DockerBox) Exec(name string, command string) error {
	defer b.client.Close()
	return b.execBox(name, command)
}

func (b *DockerBox) execBox(name string, command string) error {
	// TODO
	return errors.New("not implemented")
}

func (b *DockerBox) Open(template *model.BoxV1) error {
	defer b.client.Close()

	info, err := b.createBox(template)
	if err != nil {
		return err
	}

	return b.attachBox(info.Name, template.Shell)
}

func (b *DockerBox) attachBox(name string, command string) error {

	info, err := b.findBox(name)
	if err != nil {
		return err
	}

	containerOpts := &docker.ContainerAttachOpts{
		ContainerId: info.Id,
		Shell:       command,
		InStream:    b.opts.Streams.In,
		OutStream:   b.opts.Streams.Out,
		ErrStream:   b.opts.Streams.Err,
		IsTty:       b.opts.Streams.IsTty,
		OnContainerAttachCallback: func() {
			b.opts.EventBus.Publish(newContainerAttachDockerLoaderEvent())
		},
		OnStreamCloseCallback: func() {
			b.opts.EventBus.Publish(newContainerAttachExitDockerEvent(info.Id))
		},
		OnStreamErrorCallback: func(err error) {
			b.opts.EventBus.Publish(newContainerAttachErrorDockerEvent(info.Id, err))
		},
	}
	b.opts.EventBus.Publish(newContainerAttachDockerEvent(info.Id, info.Name, command))
	return b.client.ContainerAttach(containerOpts)
}

func (b *DockerBox) Copy(name string, from string, to string) error {
	defer b.client.Close()
	// TODO
	return errors.New("not implemented")
}

func (b *DockerBox) List() ([]box.BoxInfo, error) {
	defer b.client.Close()
	return b.listBoxes()
}

func (b *DockerBox) listBoxes() ([]box.BoxInfo, error) {

	containers, err := b.client.ContainerList(model.BoxPrefixName)
	if err != nil {
		return nil, err
	}
	var result []box.BoxInfo
	for index, c := range containers {
		result = append(result, box.BoxInfo{Id: c.ContainerId, Name: c.ContainerName})
		b.opts.EventBus.Publish(newContainerListDockerEvent(index, c.ContainerName, c.ContainerId))
	}

	return result, nil
}

func (b *DockerBox) findBox(name string) (*box.BoxInfo, error) {

	boxes, err := b.listBoxes()
	if err != nil {
		return nil, err
	}
	for _, boxInfo := range boxes {
		if boxInfo.Name == name {
			return &boxInfo, nil
		}
	}

	return nil, errors.New("box not found")
}

func (b *DockerBox) Tunnel(string) error {
	defer b.client.Close()
	return errors.New("not supported")
}

func (b *DockerBox) Delete(name string) error {
	defer b.client.Close()

	info, err := b.findBox(name)
	if err != nil {
		return err
	}

	b.opts.EventBus.Publish(newContainerRemoveDockerEvent(info.Id))
	return b.client.ContainerRemove(info.Id)
}
