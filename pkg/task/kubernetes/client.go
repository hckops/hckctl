package kubernetes

import (
	"github.com/hckops/hckctl/pkg/client/kubernetes"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/event"
	taskModel "github.com/hckops/hckctl/pkg/task/model"
)

type KubeTaskClient struct {
	client     *kubernetes.KubeClient
	clientOpts *commonModel.KubeOptions
	eventBus   *event.EventBus
}

func NewKubeTaskClient(commonOpts *taskModel.CommonTaskOptions, kubeOpts *commonModel.KubeOptions) (*KubeTaskClient, error) {
	return newKubeTaskClient(commonOpts, kubeOpts)
}

func (task *KubeTaskClient) Provider() taskModel.TaskProvider {
	return taskModel.Kubernetes
}

func (task *KubeTaskClient) Events() *event.EventBus {
	return task.eventBus
}

func (task *KubeTaskClient) Run(opts *taskModel.RunOptions) error {
	// TODO defer task.close()
	return task.runTask(opts)
}
