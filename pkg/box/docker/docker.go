package docker

import (
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/docker"
	"github.com/hckops/hckctl/pkg/command/common"
	"github.com/hckops/hckctl/pkg/util"
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

	networkName := common.ProjectName
	networkId, err := box.client.NetworkUpsert(networkName)
	if err != nil {
		return nil, err
	}
	box.opts.EventBus.Publish(newNetworkUpsertDockerEvent(networkName, networkId))

	containerOpts := &docker.ContainerCreateOpts{
		ContainerName:    containerName,
		ContainerConfig:  containerConfig,
		HostConfig:       hostConfig,
		NetworkingConfig: buildNetworkingConfig(networkName, networkId), // all on the same network
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

func buildNetworkingConfig(networkName, networkId string) *network.NetworkingConfig {
	return &network.NetworkingConfig{EndpointsConfig: map[string]*network.EndpointSettings{networkName: {NetworkID: networkId}}}
}

func (box *DockerBox) execBox(name string, command string) error {

	info, err := box.findBox(name)
	if err != nil {
		return err
	}
	if command == model.BoxShellNone {
		return box.logsBox(info.Id)
	}

	return box.attachBox(info, command, false)
}

func (box *DockerBox) openBox(template *model.BoxV1) error {

	info, err := box.createBox(template)
	if err != nil {
		return err
	}
	if template.Shell == model.BoxShellNone {
		// stop loader
		box.opts.EventBus.Publish(newContainerAttachDockerLoaderEvent())
		return box.logsBox(info.Id)
	}

	return box.attachBox(info, template.Shell, true)
}

func (box *DockerBox) attachBox(info *model.BoxInfo, command string, removeOnExit bool) error {
	box.opts.EventBus.Publish(newContainerAttachDockerEvent(info.Id, info.Name, command))

	containerOpts := &docker.ContainerAttachOpts{
		ContainerId: info.Id,
		Shell:       command,
		InStream:    box.opts.Streams.In,
		OutStream:   box.opts.Streams.Out,
		ErrStream:   box.opts.Streams.Err,
		IsTty:       box.opts.Streams.IsTty,
		OnContainerAttachCallback: func() {
			if removeOnExit {
				// stop loader
				box.opts.EventBus.Publish(newContainerAttachDockerLoaderEvent())
			}
		},
		OnStreamCloseCallback: func() {
			box.opts.EventBus.Publish(newContainerExecExitDockerEvent(info.Id))

			if removeOnExit {
				if err := box.client.ContainerRemove(info.Id); err != nil {
					box.opts.EventBus.Publish(newContainerExecErrorDockerEvent(info.Id, errors.Wrap(err, "error container exec remove")))
				}
			}
		},
		OnStreamErrorCallback: func(err error) {
			box.opts.EventBus.Publish(newContainerExecErrorDockerEvent(info.Id, err))
		},
	}

	return box.client.ContainerAttach(containerOpts)
}

func (box *DockerBox) logsBox(containerId string) error {

	opts := &docker.ContainerLogsOpts{
		ContainerId: containerId,
		OutStream:   box.opts.Streams.Out,
		ErrStream:   box.opts.Streams.Err,
		OnStreamCloseCallback: func() {
			box.opts.EventBus.Publish(newContainerExecExitDockerEvent(containerId))
		},
		OnStreamErrorCallback: func(err error) {
			box.opts.EventBus.Publish(newContainerExecErrorDockerEvent(containerId, err))
		},
	}
	return box.client.ContainerLogs(opts)
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
