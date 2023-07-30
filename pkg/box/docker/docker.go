package docker

import (
	"fmt"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/docker"
	"github.com/hckops/hckctl/pkg/schema"
	"github.com/hckops/hckctl/pkg/util"
)

func newDockerBoxClient(commonOpts *model.CommonBoxOptions, dockerOpts *model.DockerBoxOptions) (*DockerBoxClient, error) {
	commonOpts.EventBus.Publish(newClientInitDockerEvent())

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
	box.eventBus.Publish(newClientCloseDockerEvent())
	box.eventBus.Close()
	return box.client.Close()
}

// TODO limit resources by size?
func (box *DockerBoxClient) createBox(opts *model.CreateOptions) (*model.BoxInfo, error) {

	imageName := opts.Template.ImageName()
	imagePullOpts := &docker.ImagePullOpts{
		ImageName: imageName,
		OnImagePullCallback: func() {
			box.eventBus.Publish(newImagePullDockerLoaderEvent(imageName))
		},
	}
	box.eventBus.Publish(newImagePullDockerEvent(imageName))
	if err := box.client.ImagePull(imagePullOpts); err != nil {
		// try to use existing images
		if box.clientOpts.IgnoreImagePullError {
			box.eventBus.Publish(newImagePullIgnoreDockerEvent(imageName))
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
			box.eventBus.Publish(newImageRemoveIgnoreDockerEvent(imageId, err))
		},
	}
	if err := box.client.ImageRemoveDangling(imageRemoveOpts); err != nil {
		return nil, err
	}

	// TODO add env var container override
	// TODO print environment variables

	// boxName
	containerName := opts.Template.GenerateName()
	networkPorts := opts.Template.NetworkPorts(false)
	containerConfig, err := buildContainerConfig(&containerConfigOptions{
		imageName:     opts.Template.ImageName(),
		containerName: containerName,
		ports:         networkPorts,
		labels:        opts.Labels,
	})
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
		NetworkingConfig: buildNetworkingConfig(networkName, networkId), // all on the same network
	}
	// boxId
	containerId, err := box.client.ContainerCreate(containerOpts)
	if err != nil {
		return nil, err
	}
	box.eventBus.Publish(newContainerCreateDockerEvent(opts.Template.Name, containerName, containerId))

	return &model.BoxInfo{Id: containerId, Name: containerName, Healthy: true}, nil
}

type containerConfigOptions struct {
	imageName     string
	containerName string
	ports         []model.BoxPort
	labels        map[string]string
}

func buildContainerConfig(opts *containerConfigOptions) (*container.Config, error) {

	exposedPorts := make(nat.PortSet)
	for _, port := range opts.ports {
		p, err := nat.NewPort("tcp", port.Remote)
		if err != nil {
			return nil, errors.Wrap(err, "error docker port: containerConfig")
		}
		exposedPorts[p] = struct{}{}
	}

	// TODO add Env
	return &container.Config{
		Hostname:     opts.containerName,
		Image:        opts.imageName,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		OpenStdin:    true,
		StdinOnce:    true,
		Tty:          true,
		ExposedPorts: exposedPorts,
		Labels:       opts.labels,
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

func (box *DockerBoxClient) connectBox(opts *model.ConnectOptions) error {
	if info, err := box.searchBox(opts.Name); err != nil {
		return err
	} else {
		if opts.DisableExec || opts.DisableTunnel {
			box.eventBus.Publish(newContainerExecIgnoreDockerEvent(info.Id))
		}

		return box.execBox(opts.Template, info, opts.Streams, opts.DeleteOnExit)
	}
}

func (box *DockerBoxClient) searchBox(name string) (*model.BoxInfo, error) {
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

func (box *DockerBoxClient) execBox(template *model.BoxV1, info *model.BoxInfo, streams *model.BoxStreams, deleteOnExit bool) error {
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

	if command == model.BoxShellNone {
		if deleteOnExit {
			// stop loader
			box.eventBus.Publish(newContainerExecDockerLoaderEvent())
		}
		return box.logsBox(info.Id, streams)
	}

	// TODO container inspect/describe to print the actual bound ports, not the template
	// box.publishBoxInfo(template, info)

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

func (box *DockerBoxClient) publishBoxInfo(template *model.BoxV1, info *model.BoxInfo) {
	// print open ports
	networkPorts := template.NetworkPorts(false)
	portPadding := model.PortFormatPadding(networkPorts)
	for _, port := range networkPorts {
		box.eventBus.Publish(newContainerCreatePortBindDockerEvent(info.Name, port))
		box.eventBus.Publish(newContainerCreatePortBindDockerConsoleEvent(info.Name, port, portPadding))
	}
	// TODO print environment variables
}

func (box *DockerBoxClient) logsBox(containerId string, streams *model.BoxStreams) error {
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

func (box *DockerBoxClient) describe(name string) (*model.BoxDetails, error) {
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

func toBoxDetails(container docker.ContainerDetails) (*model.BoxDetails, error) {

	labels := model.BoxLabels(container.Labels)

	size, err := labels.ToSize()
	if err != nil {
		return nil, err
	}

	var env []model.BoxEnv
	for _, e := range container.Env {
		items := strings.Split(e, "=")
		if len(items) == 2 {
			env = append(env, model.BoxEnv{
				Key:   items[0],
				Value: items[1],
			})
		}
	}

	var ports []model.BoxPort
	for _, p := range container.Ports {
		ports = append(ports, model.BoxPort{
			Alias:  model.BoxPortNone, // match with template
			Local:  p.Local,
			Remote: p.Remote,
			Public: false,
		})
	}

	return &model.BoxDetails{
		Info: newBoxInfo(container.Info),
		TemplateInfo: &model.BoxTemplateInfo{
			CachedTemplate: labels.ToCachedTemplateInfo(),
			GitTemplate:    labels.ToGitTemplateInfo(),
		},
		ProviderInfo: &model.BoxProviderInfo{
			Provider: model.Docker,
			DockerProvider: &model.DockerProviderInfo{
				Network: container.Network.Name,
				Ip:      container.Network.IpAddress,
			},
		},
		Size:    size,
		Env:     model.SortEnv(env),
		Ports:   model.SortPorts(ports),
		Created: container.Created,
	}, nil
}

func newBoxInfo(container docker.ContainerInfo) model.BoxInfo {
	return model.BoxInfo{
		Id:      container.ContainerId,
		Name:    container.ContainerName,
		Healthy: container.Healthy,
	}
}

func boxLabel() string {
	return fmt.Sprintf("%s=%s", model.LabelSchemaKind, schema.KindBoxV1.String())
}

func (box *DockerBoxClient) listBoxes() ([]model.BoxInfo, error) {

	containers, err := box.client.ContainerList(model.BoxPrefixName, boxLabel())
	if err != nil {
		return nil, err
	}

	var result []model.BoxInfo
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
