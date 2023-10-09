package docker

import (
	"github.com/hckops/hckctl/pkg/client/docker"
	commonDocker "github.com/hckops/hckctl/pkg/common/docker"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/event"
	taskModel "github.com/hckops/hckctl/pkg/task/model"
)

type DockerTaskClient struct {
	client       *docker.DockerClient
	clientOpts   *commonModel.DockerOptions
	dockerCommon *commonDocker.DockerCommonClient
	eventBus     *event.EventBus
}

func NewDockerTaskClient(commonOpts *taskModel.CommonTaskOptions, dockerOpts *commonModel.DockerOptions) (*DockerTaskClient, error) {
	return newDockerTaskClient(commonOpts, dockerOpts)
}

func (task *DockerTaskClient) Provider() taskModel.TaskProvider {
	return taskModel.Docker
}

func (task *DockerTaskClient) Events() *event.EventBus {
	return task.eventBus
}

func (task *DockerTaskClient) Run(opts *taskModel.RunOptions) error {
	defer task.close()
	return task.runTask(opts)
}
