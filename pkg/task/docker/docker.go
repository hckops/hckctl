package docker

import (
	"path"

	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/client/docker"
	commonDocker "github.com/hckops/hckctl/pkg/common/docker"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	taskModel "github.com/hckops/hckctl/pkg/task/model"
)

func newDockerTaskClient(commonOpts *taskModel.CommonTaskOptions, dockerOpts *commonModel.DockerOptions) (*DockerTaskClient, error) {

	dockerClient, err := commonDocker.NewDockerCommonClient(commonOpts.EventBus, dockerOpts)
	if err != nil {
		return nil, errors.Wrap(err, "error docker task client")
	}

	return &DockerTaskClient{
		docker:     dockerClient,
		clientOpts: dockerOpts,
		eventBus:   commonOpts.EventBus,
	}, nil
}

func (task *DockerTaskClient) close() error {
	return task.docker.Client.Close()
}

func (task *DockerTaskClient) runTask(opts *taskModel.RunOptions) error {

	// taskName
	containerName := opts.Template.GenerateName()

	// vpn sidecar
	var networkMode string
	if opts.NetworkInfo.Vpn != nil {
		if sidecarContainerId, err := task.docker.StartVpnSidecar(containerName, opts.NetworkInfo.Vpn, &docker.ContainerPortConfigOpts{}); err != nil {
			return err
		} else {
			networkMode = docker.ContainerNetworkMode(sidecarContainerId)
			// remove sidecar on exit
			defer task.docker.Client.ContainerRemove(sidecarContainerId)
		}
	} else {
		networkMode = docker.DefaultNetworkMode()
	}

	imageName := opts.Template.Image.Name()
	if err := task.docker.PullImageOffline(imageName, func() {
		task.eventBus.Publish(newImagePullDockerLoaderEvent(imageName))
	}); err != nil {
		return err
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
				HostDir:      opts.ShareDir,
				ContainerDir: commonModel.MountShareDir,
			},
		},
	})
	if err != nil {
		return err
	}

	networkName := task.clientOpts.NetworkName
	networkId, err := task.docker.Client.NetworkUpsert(networkName)
	if err != nil {
		return err
	}
	task.eventBus.Publish(newNetworkUpsertDockerEvent(networkName, networkId))
	task.eventBus.Publish(newContainerCreateDockerLoaderEvent())

	// TODO prefix log file name with unix timestamp to order files?
	logFileName := path.Join(opts.LogDir, containerName)
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
			task.docker.Client.ContainerRemove(containerId)
		},
		OnContainerCreateCallback: func(string) error { return nil },
		OnContainerWaitCallback: func(containerId string) error {
			task.eventBus.Publish(newVolumeMountDockerEvent(containerId, opts.ShareDir, commonModel.MountShareDir))

			// stop loader
			task.eventBus.Publish(newContainerWaitDockerLoaderEvent())

			// TODO prepend file with actual task/command
			// TODO add flag to TaskV1 template to use "ContainerLogsStd" if command is "help" or "version"
			// tail logs before blocking
			task.eventBus.Publish(newContainerLogDockerEvent(logFileName))
			return task.docker.Client.ContainerLogsTee(containerId, logFileName)
		},
		OnContainerStatusCallback: func(status string) {
			task.eventBus.Publish(newContainerCreateStatusDockerEvent(status))
		},
		OnContainerStartCallback: func() {},
	}
	// taskId
	containerId, err := task.docker.Client.ContainerCreate(containerOpts)
	if err != nil {
		return err
	}
	task.eventBus.Publish(newContainerCreateDockerEvent(opts.Template.Name, containerName, containerId))
	task.eventBus.Publish(newContainerLogDockerConsoleEvent(logFileName))

	// remove temporary container
	return task.docker.Client.ContainerRemove(containerId)
}
