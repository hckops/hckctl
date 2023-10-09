package kubernetes

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/client/kubernetes"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/schema"
	taskModel "github.com/hckops/hckctl/pkg/task/model"
	"github.com/hckops/hckctl/pkg/util"
)

func newKubeTaskClient(commonOpts *taskModel.CommonTaskOptions, kubeOpts *commonModel.KubeOptions) (*KubeTaskClient, error) {
	commonOpts.EventBus.Publish(newInitKubeClientEvent())

	kubeClient, err := kubernetes.NewKubeClient(kubeOpts.InCluster, kubeOpts.ConfigPath)
	if err != nil {
		return nil, errors.Wrap(err, "error kube task")
	}

	return &KubeTaskClient{
		client:     kubeClient,
		clientOpts: kubeOpts,
		eventBus:   commonOpts.EventBus,
	}, nil
}

func (task *KubeTaskClient) close() error {
	task.eventBus.Publish(newCloseKubeClientEvent())
	task.eventBus.Close()
	return task.client.Close()
}

// TODO copy vpn https://github.com/goccy/kubejob/blob/master/copy.go

func (task *KubeTaskClient) runTask(opts *taskModel.RunOptions) error {
	namespace := task.clientOpts.Namespace

	// create namespace
	if err := task.client.NamespaceApply(namespace); err != nil {
		return err
	}
	task.eventBus.Publish(newNamespaceApplyKubeEvent(namespace))

	// create job
	jobName := opts.Template.GenerateName()
	jobSpec := kubernetes.BuildJob(&kubernetes.JobOpts{
		Namespace:   namespace,
		Name:        jobName,
		Annotations: opts.Labels,
		Labels: kubernetes.BuildLabels(jobName, opts.Template.Image.Repository, opts.Template.Image.ResolveVersion(),
			map[string]string{commonModel.LabelSchemaKind: util.ToLowerKebabCase(schema.KindTaskV1.String())}),
		PodInfo: &kubernetes.PodInfo{
			Namespace:     namespace,
			PodName:       "INVALID_POD_NAME", // not used, generated suffix by kube
			ContainerName: opts.Template.Image.Repository,
			ImageName:     opts.Template.Image.Name(),
			Arguments:     opts.Arguments, // TODO verify commands and arguments
			Env:           []kubernetes.KubeEnv{},
			Resource:      &kubernetes.KubeResource{}, // TODO set default
		},
	})
	jobOpts := &kubernetes.JobCreateOpts{
		Namespace: namespace,
		Spec:      jobSpec,
		OnStatusEventCallback: func(event string) {
			task.eventBus.Publish(newJobCreateStatusKubeEvent(event))
		},
	}
	err := task.client.JobCreate(jobOpts)
	if err != nil {
		return err
	}
	task.eventBus.Publish(newJobCreateKubeEvent(namespace, jobName))

	podInfo, err := task.client.PodDescribeFromJob(jobSpec)
	if err != nil {
		return err
	}
	task.eventBus.Publish(newPodNameKubeEvent(namespace, podInfo.PodName, podInfo.ContainerName))

	// stop loader
	task.eventBus.Publish(newContainerWaitKubeLoaderEvent())

	// TODO tee
	logFileName := "TODO_LOG_FILE_NAME"

	logOpts := &kubernetes.PodLogOpts{
		Namespace: namespace,
		PodName:   podInfo.PodName,
		PodId:     podInfo.ContainerName,
	}
	task.eventBus.Publish(newPodLogKubeEvent(logFileName))
	if err := task.client.PodLog(logOpts); err != nil {
		return err
	}

	task.eventBus.Publish(newPodLogKubeConsoleEvent(logFileName))
	task.eventBus.Publish(newJobDeleteKubeEvent(namespace, jobName))
	return task.client.JobDelete(namespace, jobName)
}