package docker

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/client/docker"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	taskModel "github.com/hckops/hckctl/pkg/task/model"
)

// TODO events

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

	imageName := opts.Template.Image.Name()
	imagePullOpts := &docker.ImagePullOpts{
		ImageName: imageName,
		OnImagePullCallback: func() {
			task.eventBus.Publish(newImagePullDockerLoaderEvent(imageName))
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

	// taskName
	containerName := opts.Template.GenerateName()

	containerConfig, err := docker.BuildContainerConfig(&docker.ContainerConfigOpts{
		ImageName: imageName,
		Hostname:  "", // TODO vpn NetworkMode conflicts with Hostname containerName
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
		NetworkMode:        docker.DefaultNetworkMode(), // TODO vpn
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

	containerOpts := &docker.ContainerCreateOpts{
		ContainerName:    containerName,
		ContainerConfig:  containerConfig,
		HostConfig:       hostConfig,
		NetworkingConfig: docker.BuildNetworkingConfig(networkName, networkId), // all on the same network
		WaitStatus:       true,
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
	task.eventBus.Publish(newContainerCreateDockerLoaderEvent())

	if err := task.client.ContainerLogsStd(containerId); err != nil {
		return err
	}
	if err := task.client.ContainerRemove(containerId); err != nil {
		return err
	}
	return nil
}
