package docker

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/docker"
	"github.com/hckops/hckctl/pkg/util"
)

func newDockerBox(commonOpts *model.BoxCommonOptions, clientConfig *docker.DockerClientConfig) (*DockerBox, error) {
	commonOpts.EventBus.Publish(newClientInitDockerEvent())

	dockerClient, err := docker.NewDockerClient()
	if err != nil {
		return nil, errors.Wrap(err, "error docker box")
	}

	return &DockerBox{
		client:       dockerClient,
		clientConfig: clientConfig,
		streams:      commonOpts.Streams,
		eventBus:     commonOpts.EventBus,
	}, nil
}

func (box *DockerBox) close() error {
	box.eventBus.Publish(newClientCloseDockerEvent())
	box.eventBus.Close()
	return box.client.Close()
}

func (box *DockerBox) createBox(template *model.BoxV1) (*model.BoxInfo, error) {

	imageName := template.ImageName()
	imagePullOpts := &docker.ImagePullOpts{
		ImageName: imageName,
		OnImagePullCallback: func() {
			box.eventBus.Publish(newImagePullDockerLoaderEvent(imageName))
		},
	}
	box.eventBus.Publish(newImagePullDockerEvent(imageName))
	if err := box.client.ImagePull(imagePullOpts); err != nil {
		// TODO search existing image and move IgnoreImagePullError/AllowOffline in docker client
		// try to use existing images
		if box.clientConfig.IgnoreImagePullError {
			box.eventBus.Publish(newImagePullErrorDockerEvent(imageName))
		} else {
			return nil, err
		}
	}

	// cleanup old nightly images
	imageRemoveOpts := &docker.ImageRemoveOpts{
		OnImageRemoveCallback: func(imageId string) {
			box.eventBus.Publish(newImageRemoveDockerEvent(imageId))
		},
		OnImageRemoveErrorCallback: func(imageId string, err error) {
			box.eventBus.Publish(newImageRemoveErrorDockerEvent(imageId, err))
		},
	}
	if err := box.client.ImageRemoveDangling(imageRemoveOpts); err != nil {
		return nil, err
	}

	// TODO add env var container override
	// TODO print environment variables

	// boxName
	containerName := template.GenerateName()
	networkPorts := template.NetworkPorts(false)
	containerConfig, err := buildContainerConfig(
		template.ImageName(),
		containerName,
		networkPorts,
	)
	if err != nil {
		return nil, err
	}

	portPadding := model.PortFormatPadding(networkPorts)
	onPortBindCallback := func(port model.BoxPort) {
		box.eventBus.Publish(newContainerCreatePortBindDockerEvent(containerName, port))
		box.eventBus.Publish(newContainerCreatePortBindDockerConsoleEvent(containerName, port, portPadding))
	}
	hostConfig, err := buildHostConfig(networkPorts, onPortBindCallback)
	if err != nil {
		return nil, err
	}

	networkName := box.clientConfig.NetworkName
	networkId, err := box.client.NetworkUpsert(networkName)
	if err != nil {
		return nil, err
	}
	box.eventBus.Publish(newNetworkUpsertDockerEvent(networkName, networkId))

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
	box.eventBus.Publish(newContainerCreateDockerEvent(template.Name, containerName, containerId))

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

	// TODO add label revision: use correct revision when resolving template by name
	// TODO add label owner/managed-by: use to list instead of prefix
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
		//Labels:       map[string]string{"com.hckops.revision": "main-or-empty"},
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

// TODO common
func (box *DockerBox) connectBox(template *model.BoxV1, name string) error {
	if info, err := box.findBox(name); err != nil {
		return err
	} else {
		return box.execBox(template, info, false)
	}
}

// TODO common
func (box *DockerBox) openBox(template *model.BoxV1) error {
	if info, err := box.createBox(template); err != nil {
		return err
	} else {
		return box.execBox(template, info, true)
	}
}

// TODO common
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

func (box *DockerBox) execBox(template *model.BoxV1, info *model.BoxInfo, removeOnExit bool) error {
	command := template.Shell
	box.eventBus.Publish(newContainerExecDockerEvent(info.Id, info.Name, command))

	// TODO it should print the actual bound ports, not the template
	// box.publishBoxInfo(template, info)

	if command == model.BoxShellNone {
		if removeOnExit {
			// stop loader
			box.eventBus.Publish(newContainerExecDockerLoaderEvent())
		}
		return box.logsBox(info.Id)
	}

	containerOpts := &docker.ContainerExecOpts{
		ContainerId: info.Id,
		Shell:       command,
		InStream:    box.streams.In,
		OutStream:   box.streams.Out,
		ErrStream:   box.streams.Err,
		IsTty:       box.streams.IsTty,
		OnContainerExecCallback: func() {
			if removeOnExit {
				// stop loader
				box.eventBus.Publish(newContainerExecDockerLoaderEvent())
			}
		},
		OnStreamCloseCallback: func() {
			box.eventBus.Publish(newContainerExecExitDockerEvent(info.Id))

			if removeOnExit {
				if err := box.client.ContainerRemove(info.Id); err != nil {
					box.eventBus.Publish(newContainerExecErrorDockerEvent(info.Id, errors.Wrap(err, "error container exec remove")))
				}
			}
		},
		OnStreamErrorCallback: func(err error) {
			box.eventBus.Publish(newContainerExecErrorDockerEvent(info.Id, err))
		},
	}

	return box.client.ContainerExec(containerOpts)
}

func (box *DockerBox) publishBoxInfo(template *model.BoxV1, info *model.BoxInfo) {
	// print open ports
	networkPorts := template.NetworkPorts(false)
	portPadding := model.PortFormatPadding(networkPorts)
	for _, port := range networkPorts {
		box.eventBus.Publish(newContainerCreatePortBindDockerEvent(info.Name, port))
		box.eventBus.Publish(newContainerCreatePortBindDockerConsoleEvent(info.Name, port, portPadding))
	}
	// TODO print environment variables
}

func (box *DockerBox) logsBox(containerId string) error {
	opts := &docker.ContainerLogsOpts{
		ContainerId: containerId,
		OutStream:   box.streams.Out,
		ErrStream:   box.streams.Err,
		OnStreamCloseCallback: func() {
			box.eventBus.Publish(newContainerExecExitDockerEvent(containerId))
		},
		OnStreamErrorCallback: func(err error) {
			box.eventBus.Publish(newContainerExecErrorDockerEvent(containerId, err))
		},
	}
	return box.client.ContainerLogs(opts)
}

func (box *DockerBox) listBoxes() ([]model.BoxInfo, error) {
	// TODO list by labels (add during creation)
	containers, err := box.client.ContainerList(model.BoxPrefixName)
	if err != nil {
		return nil, err
	}
	var result []model.BoxInfo
	for index, c := range containers {
		result = append(result, model.BoxInfo{Id: c.ContainerId, Name: c.ContainerName})
		box.eventBus.Publish(newContainerListDockerEvent(index, c.ContainerName, c.ContainerId))
	}
	return result, nil
}

func (box *DockerBox) deleteBoxById(id string) error {
	box.eventBus.Publish(newContainerRemoveDockerEvent(id))
	return box.client.ContainerRemove(id)
}

func (box *DockerBox) deleteBoxByName(name string) error {
	if info, err := box.findBox(name); err != nil {
		return err
	} else {
		return box.deleteBoxById(info.Id)
	}
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
		} else {
			// silently ignore
			box.eventBus.Publish(newContainerRemoveSkippedDockerEvent(boxInfo.Id))
		}
	}
	return deleted, nil
}
