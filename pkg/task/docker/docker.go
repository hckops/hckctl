package docker

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"

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
			// remove on exit
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

	containerOpts := &docker.ContainerCreateOpts{
		ContainerName:             containerName,
		ContainerConfig:           containerConfig,
		HostConfig:                hostConfig,
		NetworkingConfig:          docker.BuildNetworkingConfig(networkName, networkId), // all on the same network
		WaitStatus:                true,                                                 // block
		OnContainerCreateCallback: func(string) error { return nil },
		OnContainerWaitCallback: func(containerId string) error {
			task.eventBus.Publish(newContainerStartDockerLoaderEvent())

			// TODO error event
			// tail logs before start waiting
			if err := task.client.ContainerLogsStd(containerId); err != nil {
				return err
			}
			return nil
		},
		OnContainerStatusCallback: func(status string) {
			task.eventBus.Publish(newContainerCreateStatusDockerEvent(status))
		},
		OnContainerStartCallback: func() {
			task.eventBus.Publish(newContainerStartDockerLoaderEvent())
		},
	}
	// taskId
	containerId, err := task.client.ContainerCreate(containerOpts)
	if err != nil {
		return err
	}
	task.eventBus.Publish(newContainerCreateDockerEvent(opts.Template.Name, containerName, containerId))

	if err := task.client.ContainerRemove(containerId); err != nil {
		return err
	}
	return nil
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
	return fmt.Sprintf("sidecar-vpn-%s", strings.Split(taskName, "-")[2])
}

func (task *DockerTaskClient) startVpnSidecar(containerName string, vpnInfo *commonModel.VpnNetworkInfo) (string, error) {

	// TODO context timeout
	imageName := "hckops/alpine-openvpn:latest"
	// base directory must exist
	vpnConfigPath := "/usr/share/client.ovpn"

	if err := task.pullImage(imageName, newVpnConnectDockerLoaderEvent(vpnInfo.Name)); err != nil {
		return "", err
	}

	containerOpts := &docker.ContainerCreateOpts{
		ContainerName:   containerName,
		ContainerConfig: docker.BuildVpnContainerConfig(imageName, vpnConfigPath),
		HostConfig:      docker.BuildVpnHostConfig(),
		WaitStatus:      false,
		OnContainerCreateCallback: func(containerId string) error {
			// TODO error event
			if err := task.client.CopyFileToContainer(containerId, vpnInfo.LocalPath, vpnConfigPath); err != nil {
				return err
			}
			return nil
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
	//task.eventBus.Publish(newContainerCreateDockerEvent(opts.Template.Name, containerName, containerId))

	return containerId, nil
}
