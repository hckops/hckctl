package task

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/event"
	"github.com/hckops/hckctl/pkg/task/docker"
	"github.com/hckops/hckctl/pkg/task/kubernetes"
	"github.com/hckops/hckctl/pkg/task/model"
)

type TaskClient interface {
	Provider() model.TaskProvider
	Events() *event.EventBus
	Run(opts *model.RunOptions) error
}

func NewTaskClient(opts *model.TaskClientOptions) (TaskClient, error) {
	commonOpts := model.NewCommonTaskOpts()
	switch opts.Provider {
	case model.Docker:
		return docker.NewDockerTaskClient(commonOpts, opts.DockerOpts)
	case model.Kubernetes:
		return kubernetes.NewKubeTaskClient(commonOpts, opts.KubeOpts)
	default:
		return nil, errors.New("invalid provider")
	}
}
