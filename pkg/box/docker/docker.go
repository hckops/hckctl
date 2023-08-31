package docker

import (
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	boxModel "github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/docker"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/schema"
)

func newDockerBoxClient(commonOpts *boxModel.CommonBoxOptions, dockerOpts *commonModel.DockerOptions) (*DockerBoxClient, error) {
	commonOpts.EventBus.Publish(newInitDockerClientEvent())

	dockerClient, err := docker.NewDockerClient()
	if err != nil {
		return nil, errors.Wrap(err, "error docker box")
	}

	return &DockerBoxClient{
		client:     dockerClient,
		clientOpts: dockerOpts,
		eventBus:   commonOpts.EventBus,
	}, nil
}

func (box *DockerBoxClient) close() error {
	box.eventBus.Publish(newCloseDockerClientEvent())
	box.eventBus.Close()
	return box.client.Close()
}

// TODO limit resources by size?
func (box *DockerBoxClient) createBox(opts *boxModel.CreateOptions) (*boxModel.BoxInfo, error) {

	imageName := opts.Template.ImageName()
	imagePullOpts := &docker.ImagePullOpts{
		ImageName: imageName,
		OnImagePullCallback: func() {
			box.eventBus.Publish(newImagePullDockerLoaderEvent(imageName))
		},
	}
	box.eventBus.Publish(newImagePullDockerEvent(imageName))
	if err := box.client.ImagePull(imagePullOpts); err != nil {
		// try to use an existing image if exists
		if box.clientOpts.IgnoreImagePullError {
			box.eventBus.Publish(newImagePullIgnoreDockerEvent(imageName))
		} else {
			// do not allow offline
			return nil, err
		}
	}

	// cleanup obsolete nightly images
	imageRemoveOpts := &docker.ImageRemoveOpts{
		OnImageRemoveCallback: func(imageId string) {
			box.eventBus.Publish(newImageRemoveDockerEvent(imageId))
		},
		OnImageRemoveErrorCallback: func(imageId string, err error) {
			// ignore error: keep images used by existing containers
			box.eventBus.Publish(newImageRemoveIgnoreDockerEvent(imageId, err))
		},
	}
	if err := box.client.ImageRemoveDangling(imageRemoveOpts); err != nil {
		return nil, err
	}

	// boxName
	containerName := opts.Template.GenerateName()

	var containerEnv []docker.ContainerEnv
	for _, e := range opts.Template.EnvironmentVariables() {
		containerEnv = append(containerEnv, docker.ContainerEnv{Key: e.Key, Value: e.Value})
	}
	networkMap := opts.Template.NetworkPorts(false)
	var containerPorts []docker.ContainerPort
	for _, p := range networkMap {
		containerPorts = append(containerPorts, docker.ContainerPort{Local: p.Local, Remote: p.Remote})
	}

	containerConfig, err := docker.BuildContainerConfig(&docker.ContainerConfigOpts{
		ImageName:     opts.Template.ImageName(),
		ContainerName: containerName,
		Env:           containerEnv,
		Ports:         containerPorts,
		Labels:        opts.Labels,
	})
	if err != nil {
		return nil, err
	}

	onPortBindCallback := func(port docker.ContainerPort) {
		box.publishPortInfo(networkMap, containerName, port)
	}
	hostConfig, err := docker.BuildHostConfig(containerPorts, onPortBindCallback)
	if err != nil {
		return nil, err
	}

	networkName := box.clientOpts.NetworkName
	networkId, err := box.client.NetworkUpsert(networkName)
	if err != nil {
		return nil, err
	}
	box.eventBus.Publish(newNetworkUpsertDockerEvent(networkName, networkId))

	containerOpts := &docker.ContainerCreateOpts{
		ContainerName:    containerName,
		ContainerConfig:  containerConfig,
		HostConfig:       hostConfig,
		NetworkingConfig: docker.BuildNetworkingConfig(networkName, networkId), // all on the same network
		OnContainerStartCallback: func() {
			for _, e := range opts.Template.EnvironmentVariables() {
				box.eventBus.Publish(newContainerCreateEnvDockerEvent(containerName, e))
				box.eventBus.Publish(newContainerCreateEnvDockerConsoleEvent(containerName, e))
			}
		},
	}
	// boxId
	containerId, err := box.client.ContainerCreate(containerOpts)
	if err != nil {
		return nil, err
	}
	box.eventBus.Publish(newContainerCreateDockerEvent(opts.Template.Name, containerName, containerId))

	return &boxModel.BoxInfo{Id: containerId, Name: containerName, Healthy: true}, nil
}

func (box *DockerBoxClient) connectBox(opts *boxModel.ConnectOptions) error {
	if info, err := box.searchBox(opts.Name); err != nil {
		return err
	} else {
		if opts.DisableExec || opts.DisableTunnel {
			box.eventBus.Publish(newContainerExecIgnoreDockerEvent(info.Id))
		}
		return box.execBox(opts.Template, info, opts.Streams, opts.DeleteOnExit)
	}
}

func (box *DockerBoxClient) searchBox(name string) (*boxModel.BoxInfo, error) {
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

func (box *DockerBoxClient) execBox(template *boxModel.BoxV1, info *boxModel.BoxInfo, streams *boxModel.BoxStreams, deleteOnExit bool) error {
	command := template.Shell
	box.eventBus.Publish(newContainerExecDockerEvent(info.Id, info.Name, command))

	restartsOpts := &docker.ContainerRestartOpts{
		ContainerId: info.Id,
		OnRestartCallback: func(status string) {
			box.eventBus.Publish(newContainerRestartDockerEvent(info.Id, status))
		},
	}
	if err := box.client.ContainerRestart(restartsOpts); err != nil {
		return err
	}

	if command == boxModel.BoxShellNone {
		if deleteOnExit {
			// stop loader
			box.eventBus.Publish(newContainerExecDockerLoaderEvent())
		}
		return box.logsBox(info.Id, streams)
	}

	// already printed for temporary box
	if !deleteOnExit {
		containerDetails, err := box.client.ContainerInspect(info.Id)
		if err != nil {
			return err
		}
		// print environment variables
		for _, e := range containerDetails.Env {
			// ignore internal variables e.g. PATH
			if _, exists := template.EnvironmentVariables()[e.Key]; exists {
				env := boxModel.BoxEnv{Key: e.Key, Value: e.Value}
				box.eventBus.Publish(newContainerCreateEnvDockerEvent(containerDetails.Info.ContainerName, env))
				box.eventBus.Publish(newContainerCreateEnvDockerConsoleEvent(containerDetails.Info.ContainerName, env))
			}
		}
		// print open ports
		for _, port := range containerDetails.Ports {
			box.publishPortInfo(template.NetworkPorts(false), containerDetails.Info.ContainerName, port)
		}
	}

	execOpts := &docker.ContainerExecOpts{
		ContainerId: info.Id,
		Shell:       command,
		InStream:    streams.In,
		OutStream:   streams.Out,
		ErrStream:   streams.Err,
		IsTty:       streams.IsTty,
		OnContainerExecCallback: func() {
			// stop loader
			box.eventBus.Publish(newContainerExecDockerLoaderEvent())
		},
		OnStreamCloseCallback: func() {
			box.eventBus.Publish(newContainerExecExitDockerEvent(info.Id))

			if deleteOnExit {
				if err := box.client.ContainerRemove(info.Id); err != nil {
					box.eventBus.Publish(newContainerExecErrorDockerEvent(info.Id, errors.Wrap(err, "error container exec remove")))
				}
			}
		},
		OnStreamErrorCallback: func(err error) {
			box.eventBus.Publish(newContainerExecErrorDockerEvent(info.Id, err))
		},
	}

	return box.client.ContainerExec(execOpts)
}

func (box *DockerBoxClient) publishPortInfo(networkMap map[string]boxModel.BoxPort, containerName string, containerPort docker.ContainerPort) {
	portPadding := boxModel.PortFormatPadding(maps.Values(networkMap))

	// actual bound port
	networkPort := networkMap[containerPort.Remote]
	networkPort.Local = containerPort.Local

	box.eventBus.Publish(newContainerCreatePortBindDockerEvent(containerName, networkPort))
	box.eventBus.Publish(newContainerCreatePortBindDockerConsoleEvent(containerName, networkPort, portPadding))
}

func (box *DockerBoxClient) logsBox(containerId string, streams *boxModel.BoxStreams) error {
	opts := &docker.ContainerLogsOpts{
		ContainerId: containerId,
		OutStream:   streams.Out,
		ErrStream:   streams.Err,
		OnStreamCloseCallback: func() {
			box.eventBus.Publish(newContainerExecExitDockerEvent(containerId))
		},
		OnStreamErrorCallback: func(err error) {
			box.eventBus.Publish(newContainerExecErrorDockerEvent(containerId, err))
		},
	}
	return box.client.ContainerLogs(opts)
}

func (box *DockerBoxClient) describeBox(name string) (*boxModel.BoxDetails, error) {
	info, err := box.searchBox(name)
	if err != nil {
		return nil, err
	}

	box.eventBus.Publish(newContainerInspectDockerEvent(info.Id))
	containerInfo, err := box.client.ContainerInspect(info.Id)
	if err != nil {
		return nil, err
	}

	return toBoxDetails(containerInfo)
}

func toBoxDetails(container docker.ContainerDetails) (*boxModel.BoxDetails, error) {

	labels := commonModel.Labels(container.Labels)

	size, err := labels.ToBoxSize()
	if err != nil {
		return nil, err
	}

	var envs []boxModel.BoxEnv
	for _, e := range container.Env {
		envs = append(envs, boxModel.BoxEnv{
			Key:   e.Key,
			Value: e.Value,
		})
	}

	var ports []boxModel.BoxPort
	for _, p := range container.Ports {
		ports = append(ports, boxModel.BoxPort{
			Alias:  boxModel.BoxPortNone, // match with template
			Local:  p.Local,
			Remote: p.Remote,
			Public: false,
		})
	}

	return &boxModel.BoxDetails{
		Info: newBoxInfo(container.Info),
		TemplateInfo: &boxModel.BoxTemplateInfo{
			CachedTemplate: labels.ToCachedTemplateInfo(),
			GitTemplate:    labels.ToGitTemplateInfo(),
		},
		ProviderInfo: &boxModel.BoxProviderInfo{
			Provider: boxModel.Docker,
			DockerProvider: &commonModel.DockerProviderInfo{
				Network: container.Network.Name,
				Ip:      container.Network.IpAddress,
			},
		},
		Size:    size,
		Env:     boxModel.SortEnv(envs),
		Ports:   boxModel.SortPorts(ports),
		Created: container.Created,
	}, nil
}

func newBoxInfo(container docker.ContainerInfo) boxModel.BoxInfo {
	return boxModel.BoxInfo{
		Id:      container.ContainerId,
		Name:    container.ContainerName,
		Healthy: container.Healthy,
	}
}

func boxLabel() string {
	return fmt.Sprintf("%s=%s", commonModel.LabelSchemaKind, schema.KindBoxV1.String())
}

func (box *DockerBoxClient) listBoxes() ([]boxModel.BoxInfo, error) {

	containers, err := box.client.ContainerList(boxModel.BoxPrefixName, boxLabel())
	if err != nil {
		return nil, err
	}

	var result []boxModel.BoxInfo
	for index, c := range containers {
		result = append(result, newBoxInfo(c))
		box.eventBus.Publish(newContainerListDockerEvent(index, c.ContainerName, c.ContainerId, c.Healthy))
	}
	return result, nil
}

func (box *DockerBoxClient) deleteBoxes(names []string) ([]string, error) {

	boxes, err := box.listBoxes()
	if err != nil {
		return nil, err
	}

	var deleted []string
	for _, boxInfo := range boxes {

		// all or filter
		if len(names) == 0 || slices.Contains(names, boxInfo.Name) {

			if err := box.client.ContainerRemove(boxInfo.Id); err == nil {
				deleted = append(deleted, boxInfo.Name)
				box.eventBus.Publish(newContainerRemoveDockerEvent(boxInfo.Id))
			} else {
				// silently ignore
				box.eventBus.Publish(newContainerRemoveIgnoreDockerEvent(boxInfo.Id))
			}
		}
	}
	return deleted, nil
}
