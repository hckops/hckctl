package docker

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/client/docker"
	commonDocker "github.com/hckops/hckctl/pkg/common/docker"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	taskModel "github.com/hckops/hckctl/pkg/task/model"
)

func newDockerTaskClient(commonOpts *taskModel.CommonTaskOptions, dockerOpts *commonModel.DockerOptions) (*DockerTaskClient, error) {

	dockerCommonClient, err := commonDocker.NewDockerCommonClient(dockerOpts, commonOpts.EventBus)
	if err != nil {
		return nil, errors.Wrap(err, "error docker task client")
	}

	return &DockerTaskClient{
		client:       dockerCommonClient.GetClient(),
		clientOpts:   dockerOpts,
		dockerCommon: dockerCommonClient,
		eventBus:     commonOpts.EventBus,
	}, nil
}

func (task *DockerTaskClient) close() error {
	return task.dockerCommon.Close()
}

func (task *DockerTaskClient) runTask(opts *taskModel.RunOptions) error {

	// pull image
	imageName := opts.Template.Image.Name()
	if err := task.dockerCommon.PullImageOffline(imageName, func() {
		task.eventBus.Publish(newImagePullDockerLoaderEvent(imageName))
	}); err != nil {
		return err
	}

	// taskName
	containerName := opts.Template.GenerateName()

	// vpn sidecar
	var networkMode string
	if opts.CommonInfo.NetworkVpn != nil {
		sidecarOpts := &commonModel.SidecarVpnInjectOpts{
			MainContainerId: containerName,
			NetworkVpn:      opts.CommonInfo.NetworkVpn,
		}
		if sidecarContainerId, err := task.dockerCommon.SidecarVpnInject(sidecarOpts, &docker.ContainerPortConfigOpts{}); err != nil {
			return err
		} else {
			networkMode = docker.ContainerNetworkMode(sidecarContainerId)
			// remove sidecar on exit
			defer task.client.ContainerRemove(sidecarContainerId)
		}
	} else {
		networkMode = docker.DefaultNetworkMode()
	}

	containerConfig, err := docker.BuildContainerConfig(&docker.ContainerConfigOpts{
		ImageName: imageName,
		Hostname:  "", // vpn NetworkMode conflicts with Hostname containerName
		Env:       []docker.ContainerEnv{},
		Ports:     []docker.ContainerPort{},
		Labels:    opts.Labels,
		Tty:       opts.StreamOpts.IsTty,
		Cmd:       opts.Arguments,
	})
	if err != nil {
		return err
	}

	hostConfig, err := docker.BuildHostConfig(&docker.ContainerHostConfigOpts{
		NetworkMode: networkMode,
		PortConfig:  &docker.ContainerPortConfigOpts{},
		Volumes: []docker.ContainerVolume{
			{
				HostDir:      opts.CommonInfo.ShareDir.LocalPath,
				ContainerDir: opts.CommonInfo.ShareDir.RemotePath,
			},
		},
	})
	if err != nil {
		return err
	}

	networkName := task.clientOpts.NetworkName
	networkId, err := task.client.NetworkUpsert(networkName)
	if err != nil {
		return err
	}
	task.eventBus.Publish(newNetworkUpsertDockerEvent(networkName, networkId))
	task.eventBus.Publish(newContainerCreateDockerLoaderEvent())

	logFileName := opts.GenerateLogFileName(taskModel.Docker, containerName)
	containerOpts := &docker.ContainerCreateOpts{
		ContainerName:    containerName,
		ContainerConfig:  containerConfig,
		HostConfig:       hostConfig,
		NetworkingConfig: docker.BuildNetworkingConfig(networkName, networkId), // all on the same network
		WaitStatus:       true,                                                 // block
		CaptureInterrupt: true,
		OnContainerInterruptCallback: func(containerId string) {
			// returns control to runTask, it will correctly invoke defer to remove the sidecar
			// unless it's interrupted while the sidecar is being created
			task.eventBus.Publish(newContainerRemoveDockerEvent(containerId))
			task.client.ContainerRemove(containerId)
		},
		OnContainerCreateCallback: func(string) error { return nil },
		OnContainerWaitCallback: func(containerId string) error {
			task.eventBus.Publish(newVolumeMountDockerEvent(containerId, opts.CommonInfo.ShareDir.LocalPath, opts.CommonInfo.ShareDir.RemotePath))

			// stop loader
			task.eventBus.Publish(newContainerWaitDockerLoaderEvent())

			// TODO prepend file content with actual task/command
			// TODO add flag to TaskV1 template to use "ContainerLogsStd" if command is "help" or "version"
			// tail logs before blocking
			task.eventBus.Publish(newContainerLogDockerEvent(logFileName))
			return task.client.ContainerLogsTee(containerId, logFileName)
		},
		OnContainerStatusCallback: func(status string) {
			task.eventBus.Publish(newContainerCreateStatusDockerEvent(status))
		},
		OnContainerStartCallback: func() {},
	}
	// taskId
	containerId, err := task.client.ContainerCreate(containerOpts)
	if err != nil {
		return err
	}
	task.eventBus.Publish(newContainerCreateDockerEvent(opts.Template.Name, containerName, containerId))
	task.eventBus.Publish(newContainerLogDockerConsoleEvent(logFileName))

	// remove temporary container
	task.eventBus.Publish(newContainerRemoveDockerEvent(containerId))
	return task.client.ContainerRemove(containerId)
}
