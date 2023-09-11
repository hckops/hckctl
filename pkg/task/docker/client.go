package docker

import (
	"github.com/hckops/hckctl/pkg/client/docker"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/event"
	taskModel "github.com/hckops/hckctl/pkg/task/model"
)

type DockerTaskClient struct {
	client     *docker.DockerClient
	clientOpts *commonModel.DockerOptions
	eventBus   *event.EventBus
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
	return task.runTask(opts)
}
