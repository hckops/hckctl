package docker

import (
	"fmt"
	"strings"

	"github.com/docker/docker/api/types/container" // TODO remove
	"github.com/docker/docker/api/types/network"   // TODO remove
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
func (box *DockerBoxClient) createBox(opts *model.TemplateOptions) (*model.BoxInfo, error) {

	imageName := opts.Template.ImageName()
	imagePullOpts := &docker.ImagePullOpts{
		ImageName: imageName,
		OnImagePullCallback: func() {
			box.eventBus.Publish(newImagePullDockerLoaderEvent(imageName))
		},
	}
	box.eventBus.Publish(newImagePullDockerEvent(imageName))
	if err := box.client.ImagePull(imagePullOpts); err != nil {
		// TODO search existing image
		// try to use existing images
		if box.clientOpts.IgnoreImagePullError {
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

// TODO refactor in docker client
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

// TODO refactor in docker client
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

// TODO refactor in docker client
func buildNetworkingConfig(networkName, networkId string) *network.NetworkingConfig {
	return &network.NetworkingConfig{EndpointsConfig: map[string]*network.EndpointSettings{networkName: {NetworkID: networkId}}}
}

// TODO common
func (box *DockerBoxClient) connectBox(template *model.BoxV1, tunnelOpts *model.TunnelOptions, name string) error {
	if info, err := box.findBox(name); err != nil {
		return err
	} else {
		return box.execBox(template, info, tunnelOpts, false)
	}
}

// TODO common
func (box *DockerBoxClient) openBox(templateOpts *model.TemplateOptions, tunnelOpts *model.TunnelOptions) error {
	if info, err := box.createBox(templateOpts); err != nil {
		return err
	} else {
		return box.execBox(templateOpts.Template, info, tunnelOpts, true)
	}
}

// TODO common
func (box *DockerBoxClient) findBox(name string) (*model.BoxInfo, error) {
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

func (box *DockerBoxClient) execBox(template *model.BoxV1, info *model.BoxInfo, tunnelOpts *model.TunnelOptions, removeOnExit bool) error {
	command := template.Shell
	box.eventBus.Publish(newContainerExecDockerEvent(info.Id, info.Name, command))

	// TODO if BoxInfo not Healthy attempt restart

	// TODO see command ValidateTunnelFlag ?!
	// TODO TunnelOnly > skip exec
	// TODO NoTunnel > print console warning: flag ignored

	// TODO it should print the actual bound ports, not the template
	// box.publishBoxInfo(template, info)

	if command == model.BoxShellNone {
		if removeOnExit {
			// stop loader
			box.eventBus.Publish(newContainerExecDockerLoaderEvent())
		}
		return box.logsBox(info.Id, tunnelOpts)
	}

	containerOpts := &docker.ContainerExecOpts{
		ContainerId: info.Id,
		Shell:       command,
		InStream:    tunnelOpts.Streams.In,
		OutStream:   tunnelOpts.Streams.Out,
		ErrStream:   tunnelOpts.Streams.Err,
		IsTty:       tunnelOpts.Streams.IsTty,
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

func (box *DockerBoxClient) logsBox(containerId string, tunnelOpts *model.TunnelOptions) error {
	opts := &docker.ContainerLogsOpts{
		ContainerId: containerId,
		OutStream:   tunnelOpts.Streams.Out,
		ErrStream:   tunnelOpts.Streams.Err,
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
	info, err := box.findBox(name)
	if err != nil {
		return nil, err
	}

	containerInfo, err := box.client.ContainerInspect(info.Id)
	if err != nil {
		return nil, err
	}

	return toBoxDetails(containerInfo)
}

// TODO test
func toBoxDetails(container docker.ContainerDetails) (*model.BoxDetails, error) {

	labels := model.BoxLabels(container.Labels)

	size, err := labels.ToSize()
	if err != nil {
		return nil, err
	}

	// TODO filter by prefix e.g. "HCK_"
	var env []model.BoxEnv
	for _, e := range container.Env {
		items := strings.Split(e, "=")
		env = append(env, model.BoxEnv{
			Key:   items[0],
			Value: items[1],
		})
	}

	var ports []model.BoxPort
	for _, p := range container.Ports {
		ports = append(ports, model.BoxPort{
			Alias:  "TODO", // TODO match with template
			Local:  p.Local,
			Remote: p.Remote,
			Public: false,
		})
	}

	return &model.BoxDetails{
		Info:     newBoxInfo(container.Info),
		Provider: model.Docker,
		Size:     size,
		TemplateInfo: &model.BoxTemplateInfo{
			LocalTemplate: labels.ToLocalTemplateInfo(),
			GitTemplate:   labels.ToGitTemplateInfo(),
		},
		Env:   env,
		Ports: ports,
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

		if len(names) == 0 || slices.Contains(names, boxInfo.Name) {

			if err := box.client.ContainerRemove(boxInfo.Id); err == nil {
				deleted = append(deleted, boxInfo.Name)
				box.eventBus.Publish(newContainerRemoveDockerEvent(boxInfo.Id))
			} else {
				// silently ignore
				box.eventBus.Publish(newContainerRemoveSkippedDockerEvent(boxInfo.Id))
			}
		}
	}
	return deleted, nil
}
