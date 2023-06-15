package docker

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/docker"
	"github.com/hckops/hckctl/pkg/util"
	"github.com/pkg/errors"
	"strings"
)

func newDockerBox(opts *model.BoxOpts) (*DockerBox, error) {
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

func (box *DockerBox) close() error {
	box.opts.EventBus.Publish(newClientCloseDockerEvent())
	box.opts.EventBus.Close()
	return box.client.Close()
}

func (box *DockerBox) createBox(template *model.BoxV1) (*model.BoxInfo, error) {

	imageName := template.ImageName()
	imagePullOpts := &docker.ImagePullOpts{
		ImageName: imageName,
		OnImagePullCallback: func() {
			box.opts.EventBus.Publish(newImagePullDockerLoaderEvent(imageName))
		},
	}
	box.opts.EventBus.Publish(newImagePullDockerEvent(imageName))
	if err := box.client.ImagePull(imagePullOpts); err != nil {
		return nil, err
	}

	// cleanup old nightly images
	imageRemoveOpts := &docker.ImageRemoveOpts{
		OnImageRemoveCallback: func(imageId string) {
			box.opts.EventBus.Publish(newImageRemoveDockerEvent(imageId))
		},
		OnImageRemoveErrorCallback: func(imageId string, err error) {
			box.opts.EventBus.Publish(newImageRemoveErrorDockerEvent(imageId, err))
		},
	}
	if err := box.client.ImageRemoveDangling(imageRemoveOpts); err != nil {
		return nil, err
	}

	// boxName
	containerName := template.GenerateName()
	// skip not supported virtual-* ports
	var networkPorts []model.BoxPort
	for _, networkPort := range template.NetworkPorts() {
		if strings.HasPrefix(networkPort.Alias, model.BoxPrefixVirtualPort) {
			box.opts.EventBus.Publish(newContainerCreateSkipVirtualPortDockerEvent(containerName, networkPort))
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
		box.opts.EventBus.Publish(newContainerCreatePortBindDockerEvent(containerName, port))
		box.opts.EventBus.Publish(newContainerCreatePortBindDockerConsoleEvent(containerName, port))
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
	containerId, err := box.client.ContainerCreate(containerOpts)
	if err != nil {
		return nil, err
	}
	box.opts.EventBus.Publish(newContainerCreateDockerEvent(template.Name, containerName, containerId))

	return &model.BoxInfo{Id: containerId, Name: containerName}, nil
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
		//Labels:       map[string]string{"com.hckops.revision": "main-or-empty"}, // TODO use correct revision when resolving template by name
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

func (box *DockerBox) execBox(name string, command string) error {
	// TODO
	return errors.New("not implemented")
}

func (box *DockerBox) openBox(template *model.BoxV1) error {
	info, err := box.createBox(template)
	if err != nil {
		return err
	}
	return box.attachAndRemoveBox(info, template.Shell)
}

func (box *DockerBox) attachAndRemoveBox(info *model.BoxInfo, command string) error {
	box.opts.EventBus.Publish(newContainerAttachDockerEvent(info.Id, info.Name, command))

	containerOpts := &docker.ContainerAttachOpts{
		ContainerId: info.Id,
		Shell:       command,
		InStream:    box.opts.Streams.In,
		OutStream:   box.opts.Streams.Out,
		ErrStream:   box.opts.Streams.Err,
		IsTty:       box.opts.Streams.IsTty,
		OnContainerAttachCallback: func() {
			box.opts.EventBus.Publish(newContainerAttachDockerLoaderEvent())
		},
		OnStreamCloseCallback: func() {
			box.opts.EventBus.Publish(newContainerAttachExitDockerEvent(info.Id))
		},
		OnStreamErrorCallback: func(err error) {
			box.opts.EventBus.Publish(newContainerAttachErrorDockerEvent(info.Id, err))
		},
	}

	return box.client.ContainerAttachAndRemove(containerOpts)
}

func (box *DockerBox) listBoxes() ([]model.BoxInfo, error) {

	containers, err := box.client.ContainerList(model.BoxPrefixName)
	if err != nil {
		return nil, err
	}
	var result []model.BoxInfo
	for index, c := range containers {
		result = append(result, model.BoxInfo{Id: c.ContainerId, Name: c.ContainerName})
		box.opts.EventBus.Publish(newContainerListDockerEvent(index, c.ContainerName, c.ContainerId))
	}

	return result, nil
}

func (box *DockerBox) findBox(name string) (*model.BoxInfo, error) {
	boxes, err := box.listBoxes()
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

func (box *DockerBox) deleteBoxById(id string) error {
	box.opts.EventBus.Publish(newContainerRemoveDockerEvent(id))
	return box.client.ContainerRemove(id)
}

func (box *DockerBox) deleteBoxByName(name string) error {
	info, err := box.findBox(name)
	if err != nil {
		return err
	}
	return box.deleteBoxById(info.Id)
}

func (box *DockerBox) deleteBoxes() ([]model.BoxInfo, error) {
	boxes, err := box.listBoxes()
	if err != nil {
		return nil, err
	}
	var deleted []model.BoxInfo
	for _, boxInfo := range boxes {
		if err := box.deleteBoxById(boxInfo.Id); err == nil {
			deleted = append(deleted, boxInfo)
		}
	}
	return deleted, nil
}
