package kubernetes

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/client/kubernetes"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	taskModel "github.com/hckops/hckctl/pkg/task/model"
)

func newKubeTaskClient(commonOpts *taskModel.CommonTaskOptions, kubeOpts *commonModel.KubeOptions) (*KubeTaskClient, error) {
	// TODO commonOpts.EventBus.Publish(newInitKubeClientEvent())

	kubeClient, err := kubernetes.NewKubeClient(kubeOpts.InCluster, kubeOpts.ConfigPath)
	if err != nil {
		return nil, errors.Wrap(err, "error kube box")
	}

	return &KubeTaskClient{
		client:     kubeClient,
		clientOpts: kubeOpts,
		eventBus:   commonOpts.EventBus,
	}, nil
}

func (task *KubeTaskClient) runTask(opts *taskModel.RunOptions) error {
	return errors.New("run task not implemented")
}
