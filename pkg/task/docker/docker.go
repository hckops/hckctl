package docker

import (
	"fmt"
	"path"
	"strings"

	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/client/docker"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	taskModel "github.com/hckops/hckctl/pkg/task/model"
)

func newDockerTaskClient(commonOpts *taskModel.CommonTaskOptions, dockerOpts *commonModel.DockerOptions) (*DockerTaskClient, error) {
	commonOpts.EventBus.Publish(newInitDockerClientEvent())

	dockerClient, err := docker.NewDockerClient()
	if err != nil {
		return nil, errors.Wrap(err, "error docker task")
	}

	return &DockerTaskClient{
		client:     dockerClient,
		clientOpts: dockerOpts,
		eventBus:   commonOpts.EventBus,
	}, nil
}

func (task *DockerTaskClient) close() error {
	task.eventBus.Publish(newCloseDockerClientEvent())
	task.eventBus.Close()
	return task.client.Close()
}

func (task *DockerTaskClient) runTask(opts *taskModel.RunOptions) error {

	// taskName
	containerName := opts.Template.GenerateName()

	// vpn sidecar
	var networkMode string
	if opts.NetworkInfo.Vpn != nil {
		sidecarContainerName := buildVpnSidecarName(containerName)
		if sidecarContainerId, err := task.startVpnSidecar(sidecarContainerName, opts.NetworkInfo.Vpn); err != nil {
			return err
		} else {
			networkMode = docker.ContainerNetworkMode(sidecarContainerId)
			// remove sidecar on exit
			defer task.client.ContainerRemove(sidecarContainerId)
		}
	} else {
		networkMode = docker.DefaultNetworkMode()
	}

	imageName := opts.Template.Image.Name()
	if err := task.pullImage(imageName, newImagePullDockerLoaderEvent(imageName)); err != nil {
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
		NetworkMode:        networkMode,
		Ports:              []docker.ContainerPort{},
		OnPortBindCallback: func(docker.ContainerPort) {},
		Volumes: []docker.ContainerVolume{
			{
				HostDir:      opts.ShareDir,
				ContainerDir: commonModel.MountedShareDir,
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

	// TODO prefix file name with unix timestamp to order file?
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
			task.client.ContainerRemove(containerId)
		},
		OnContainerCreateCallback: func(string) error { return nil },
		OnContainerWaitCallback: func(containerId string) error {
			task.eventBus.Publish(newVolumeMountDockerEvent(containerId, opts.ShareDir, commonModel.MountedShareDir))

			// stop loader
			task.eventBus.Publish(newContainerWaitDockerLoaderEvent())

			// TODO prepend file with actual task/command
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
	return task.client.ContainerRemove(containerId)
}

func (task *DockerTaskClient) pullImage(imageName string, pullEvent *dockerTaskEvent) error {
	imagePullOpts := &docker.ImagePullOpts{
		ImageName: imageName,
		OnImagePullCallback: func() {
			task.eventBus.Publish(pullEvent)
		},
	}
	task.eventBus.Publish(newImagePullDockerEvent(imageName))
	if err := task.client.ImagePull(imagePullOpts); err != nil {
		// try to use an existing image if exists
		if task.clientOpts.IgnoreImagePullError {
			task.eventBus.Publish(newImagePullIgnoreDockerEvent(imageName))
		} else {
			// do not allow offline
			return err
		}
	}

	// cleanup obsolete nightly images
	imageRemoveOpts := &docker.ImageRemoveOpts{
		OnImageRemoveCallback: func(imageId string) {
			task.eventBus.Publish(newImageRemoveDockerEvent(imageId))
		},
		OnImageRemoveErrorCallback: func(imageId string, err error) {
			// ignore error: keep images used by existing containers
			task.eventBus.Publish(newImageRemoveIgnoreDockerEvent(imageId, err))
		},
	}
	if err := task.client.ImageRemoveDangling(imageRemoveOpts); err != nil {
		return err
	}

	return nil
}

func buildVpnSidecarName(taskName string) string {
	tokens := strings.Split(taskName, "-")
	return fmt.Sprintf("sidecar-vpn-%s", tokens[len(tokens)-1])
}

func (task *DockerTaskClient) startVpnSidecar(containerName string, vpnInfo *commonModel.VpnNetworkInfo) (string, error) {

	imageName := commonModel.SidecarVpnImageName
	// base directory "/usr/share" must exist
	vpnConfigPath := "/usr/share/client.ovpn"

	if err := task.pullImage(imageName, newVpnConnectDockerLoaderEvent(vpnInfo.Name)); err != nil {
		return "", err
	}

	containerOpts := &docker.ContainerCreateOpts{
		ContainerName:    containerName,
		ContainerConfig:  docker.BuildVpnContainerConfig(imageName, vpnConfigPath),
		HostConfig:       docker.BuildVpnHostConfig(),
		WaitStatus:       false,
		CaptureInterrupt: false, // edge case: killing this while creating will produce an orphan sidecar container
		OnContainerCreateCallback: func(containerId string) error {
			// upload openvpn config file
			return task.client.CopyFileToContainer(containerId, vpnInfo.LocalPath, vpnConfigPath)
		},
		OnContainerStatusCallback: func(status string) {
			task.eventBus.Publish(newContainerCreateStatusDockerEvent(status))
		},
		OnContainerStartCallback: func() {},
	}
	// sidecarId
	containerId, err := task.client.ContainerCreate(containerOpts)
	if err != nil {
		return "", err
	}
	task.eventBus.Publish(newContainerCreateDockerEvent("sidecar-vpn", containerName, containerId))

	return containerId, nil
}
