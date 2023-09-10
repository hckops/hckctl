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
	// TODO
	return &DockerTaskClient{}, nil
}

func (lab *DockerTaskClient) Provider() taskModel.TaskProvider {
	return taskModel.Docker
}

func (lab *DockerTaskClient) Events() *event.EventBus {
	return lab.eventBus
}

func (lab *DockerTaskClient) Run(opts *taskModel.CreateOptions) error {
	// TODO
	return nil
}
