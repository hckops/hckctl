package docker

import (
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	boxModel "github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/docker"
	commonDocker "github.com/hckops/hckctl/pkg/common/docker"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/schema"
)

func newDockerBoxClient(commonOpts *boxModel.CommonBoxOptions, dockerOpts *commonModel.DockerOptions) (*DockerBoxClient, error) {

	dockerClient, err := commonDocker.NewDockerCommonClient(commonOpts.EventBus, dockerOpts)
	if err != nil {
		return nil, errors.Wrap(err, "error docker box client")
	}

	return &DockerBoxClient{
		docker:     dockerClient,
		clientOpts: dockerOpts,
		eventBus:   commonOpts.EventBus,
	}, nil
}

func (box *DockerBoxClient) close() error {
	return box.docker.Client.Close()
}

// TODO limit resources by size?
func (box *DockerBoxClient) createBox(opts *boxModel.CreateOptions) (*boxModel.BoxInfo, error) {

	// pull image
	imageName := opts.Template.Image.Name()
	if err := box.docker.PullImageOffline(imageName, func() {
		box.eventBus.Publish(newImagePullDockerLoaderEvent(imageName))
	}); err != nil {
		return nil, err
	}

	// boxName
	containerName := opts.Template.GenerateName()

	// ports
	networkMap := opts.Template.NetworkPorts(false)
	var containerPorts []docker.ContainerPort
	for _, p := range networkMap {
		containerPorts = append(containerPorts, docker.ContainerPort{Local: p.Local, Remote: p.Remote})
	}
	portConfig := &docker.ContainerPortConfigOpts{
		Ports: containerPorts,
		OnPortBindCallback: func(port docker.ContainerPort) {
			box.publishPortInfo(networkMap, containerName, port)
		},
	}

	// vpn sidecar
	var hostname string
	var networkMode string
	if opts.NetworkInfo.Vpn != nil {
		// set all networks configs on the sidecar to avoid options conflicts
		if sidecarContainerId, err := box.docker.StartVpnSidecar(containerName, opts.NetworkInfo.Vpn, portConfig); err != nil {
			return nil, err
		} else {
			// fix conflicting options: hostname and the network mode
			hostname = ""

			// fix conflicting options: port exposing and the container type network mode
			containerPorts = []docker.ContainerPort{}

			// fix conflicting options: port publishing and the container type network mode
			portConfig = &docker.ContainerPortConfigOpts{}

			// use vpn network
			networkMode = docker.ContainerNetworkMode(sidecarContainerId)
		}
	} else {
		// defaults
		hostname = containerName
		networkMode = docker.DefaultNetworkMode()
	}

	var containerEnv []docker.ContainerEnv
	for _, e := range opts.Template.EnvironmentVariables() {
		containerEnv = append(containerEnv, docker.ContainerEnv{Key: e.Key, Value: e.Value})
	}

	containerConfig, err := docker.BuildContainerConfig(&docker.ContainerConfigOpts{
		ImageName: imageName,
		Hostname:  hostname,
		Env:       containerEnv,
		Ports:     containerPorts,
		Labels:    opts.Labels,
		Tty:       true, // always
		Cmd:       []string{},
	})
	if err != nil {
		return nil, err
	}

	hostConfig, err := docker.BuildHostConfig(&docker.ContainerHostConfigOpts{
		NetworkMode: networkMode,
		PortConfig:  portConfig,
		Volumes: []docker.ContainerVolume{
			{
				HostDir:      opts.ShareDir,
				ContainerDir: commonModel.MountShareDir,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	networkName := box.clientOpts.NetworkName
	networkId, err := box.docker.Client.NetworkUpsert(networkName)
	if err != nil {
		return nil, err
	}
	box.eventBus.Publish(newNetworkUpsertDockerEvent(networkName, networkId))

	containerOpts := &docker.ContainerCreateOpts{
		ContainerName:                containerName,
		ContainerConfig:              containerConfig,
		HostConfig:                   hostConfig,
		NetworkingConfig:             docker.BuildNetworkingConfig(networkName, networkId), // all on the same network
		WaitStatus:                   false,
		CaptureInterrupt:             false,           // TODO ???
		OnContainerInterruptCallback: func(string) {}, // TODO ???
		OnContainerCreateCallback:    func(string) error { return nil },
		OnContainerWaitCallback:      func(string) error { return nil },
		OnContainerStatusCallback: func(status string) {
			box.eventBus.Publish(newContainerCreateStatusDockerEvent(status))
		},
		OnContainerStartCallback: func() {
			for _, e := range opts.Template.EnvironmentVariables() {
				box.eventBus.Publish(newContainerCreateEnvDockerEvent(containerName, e))
				box.eventBus.Publish(newContainerCreateEnvDockerConsoleEvent(containerName, e))
			}
		},
	}
	// boxId
	containerId, err := box.docker.Client.ContainerCreate(containerOpts)
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
		return box.execBox(opts.Template, info, opts.StreamOpts, opts.DeleteOnExit)
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

func (box *DockerBoxClient) execBox(template *boxModel.BoxV1, info *boxModel.BoxInfo, streamOpts *commonModel.StreamOptions, deleteOnExit bool) error {
	command := template.Shell
	box.eventBus.Publish(newContainerExecDockerEvent(info.Id, info.Name, command))

	// attempt to restart all associated sidecars
	sidecars, err := box.docker.GetSidecars(info.Name)
	if err != nil {
		return err
	}
	for _, sidecar := range sidecars {
		restartsOpts := &docker.ContainerRestartOpts{
			ContainerId: sidecar.Id,
			OnRestartCallback: func(status string) {
				box.eventBus.Publish(newContainerRestartDockerEvent(sidecar.Id, status))
			},
		}
		if err := box.docker.Client.ContainerRestart(restartsOpts); err != nil {
			return err
		}
	}

	restartsOpts := &docker.ContainerRestartOpts{
		ContainerId: info.Id,
		OnRestartCallback: func(status string) {
			box.eventBus.Publish(newContainerRestartDockerEvent(info.Id, status))
		},
	}
	if err := box.docker.Client.ContainerRestart(restartsOpts); err != nil {
		return err
	}

	if command == boxModel.BoxShellNone {
		if deleteOnExit {
			// stop loader
			box.eventBus.Publish(newContainerExecDockerLoaderEvent())
		}
		return box.logsBox(info.Id, streamOpts)
	}

	// already printed for temporary box
	if !deleteOnExit {
		containerDetails, err := box.docker.Client.ContainerInspect(info.Id)
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
		// print sidecar ports
		for _, sidecar := range sidecars {
			sidecarDetails, err := box.docker.Client.ContainerInspect(sidecar.Id)
			if err != nil {
				return err
			}
			for _, port := range sidecarDetails.Ports {
				box.publishPortInfo(template.NetworkPorts(false), containerDetails.Info.ContainerName, port)
			}
		}
	}

	execOpts := &docker.ContainerExecOpts{
		ContainerId: info.Id,
		Shell:       command,
		InStream:    streamOpts.In,
		OutStream:   streamOpts.Out,
		ErrStream:   streamOpts.Err,
		IsTty:       streamOpts.IsTty,
		OnContainerExecCallback: func() {
			// stop loader
			box.eventBus.Publish(newContainerExecDockerLoaderEvent())
		},
		OnStreamCloseCallback: func() {
			box.eventBus.Publish(newContainerExecExitDockerEvent(info.Id))

			if deleteOnExit {
				for _, sidecar := range sidecars {
					if err := box.docker.Client.ContainerRemove(sidecar.Id); err != nil {
						box.eventBus.Publish(newContainerExecErrorDockerEvent(sidecar.Id, errors.Wrap(err, "error sidecar exec remove")))
					}
				}

				if err := box.docker.Client.ContainerRemove(info.Id); err != nil {
					box.eventBus.Publish(newContainerExecErrorDockerEvent(info.Id, errors.Wrap(err, "error container exec remove")))
				}
			}
		},
		OnStreamErrorCallback: func(err error) {
			box.eventBus.Publish(newContainerExecErrorDockerEvent(info.Id, err))
		},
	}

	return box.docker.Client.ContainerExec(execOpts)
}

func (box *DockerBoxClient) publishPortInfo(networkMap map[string]boxModel.BoxPort, containerName string, containerPort docker.ContainerPort) {
	portPadding := boxModel.PortFormatPadding(maps.Values(networkMap))

	// actual bound port
	networkPort := networkMap[containerPort.Remote]
	networkPort.Local = containerPort.Local

	box.eventBus.Publish(newContainerCreatePortBindDockerEvent(containerName, networkPort))
	box.eventBus.Publish(newContainerCreatePortBindDockerConsoleEvent(containerName, networkPort, portPadding))
}

func (box *DockerBoxClient) logsBox(containerId string, streamOpts *commonModel.StreamOptions) error {
	opts := &docker.ContainerLogsOpts{
		ContainerId: containerId,
		OutStream:   streamOpts.Out,
		ErrStream:   streamOpts.Err,
		OnStreamCloseCallback: func() {
			box.eventBus.Publish(newContainerExecExitDockerEvent(containerId))
		},
		OnStreamErrorCallback: func(err error) {
			box.eventBus.Publish(newContainerExecErrorDockerEvent(containerId, err))
		},
	}
	return box.docker.Client.ContainerLogs(opts)
}

func (box *DockerBoxClient) describeBox(name string) (*boxModel.BoxDetails, error) {
	info, err := box.searchBox(name)
	if err != nil {
		return nil, err
	}

	box.eventBus.Publish(newContainerInspectDockerEvent(info.Id))
	containerInfo, err := box.docker.Client.ContainerInspect(info.Id)
	if err != nil {
		return nil, err
	}

	return toBoxDetails(containerInfo)
}

func toBoxDetails(container docker.ContainerDetails) (*boxModel.BoxDetails, error) {

	labels := commonModel.Labels(container.Labels)

	size, err := boxModel.ToBoxSize(labels)
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

	containers, err := box.docker.Client.ContainerList(boxModel.BoxPrefixName, boxLabel())
	if err != nil {
		return nil, err
	}

	var boxes []boxModel.BoxInfo
	for index, c := range containers {
		boxes = append(boxes, newBoxInfo(c))
		box.eventBus.Publish(newContainerListDockerEvent(index, c.ContainerName, c.ContainerId, c.Healthy))
	}
	return boxes, nil
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

			if err := box.docker.Client.ContainerRemove(boxInfo.Id); err == nil {
				deleted = append(deleted, boxInfo.Name)
				box.eventBus.Publish(newContainerRemoveDockerEvent(boxInfo.Id))
			} else {
				// silently ignore
				box.eventBus.Publish(newContainerRemoveIgnoreDockerEvent(boxInfo.Id))
			}

			// delete all sidecars
			sidecars, _ := box.docker.GetSidecars(boxInfo.Name)
			for _, sidecar := range sidecars {
				if err := box.docker.Client.ContainerRemove(sidecar.Id); err != nil {
					box.eventBus.Publish(newContainerRemoveDockerEvent(sidecar.Id))
				} else {
					// silently ignore
					box.eventBus.Publish(newContainerRemoveIgnoreDockerEvent(sidecar.Id))
				}
			}
		}
	}
	return deleted, nil
}
