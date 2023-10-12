package kubernetes

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/client/kubernetes"
	commonKube "github.com/hckops/hckctl/pkg/common/kubernetes"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/schema"
	taskModel "github.com/hckops/hckctl/pkg/task/model"
	"github.com/hckops/hckctl/pkg/util"
)

func newKubeTaskClient(commonOpts *taskModel.CommonTaskOptions, kubeOpts *commonModel.KubeOptions) (*KubeTaskClient, error) {

	kubeCommonClient, err := commonKube.NewKubeCommonClient(kubeOpts, commonOpts.EventBus)
	if err != nil {
		return nil, errors.Wrap(err, "error kube task client")
	}

	return &KubeTaskClient{
		client:     kubeCommonClient.GetClient(),
		clientOpts: kubeOpts,
		kubeCommon: kubeCommonClient,
		eventBus:   commonOpts.EventBus,
	}, nil
}

func (task *KubeTaskClient) close() error {
	return task.kubeCommon.Close()
}

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
			Arguments:     opts.Arguments,
			Env:           []kubernetes.KubeEnv{},
			// TODO set default
			Resource: &kubernetes.KubeResource{
				Memory: "1024Mi",
				Cpu:    "1000m",
			},
		},
	})

	// create secret and inject sidecar-vpn
	if opts.NetworkInfo.Vpn != nil {
		sidecarOpts := &commonModel.SidecarVpnInjectOpts{
			MainContainerName: jobName,
			VpnInfo:           opts.NetworkInfo.Vpn,
		}
		if err := task.kubeCommon.SidecarVpnInject(namespace, sidecarOpts, &jobSpec.Spec.Template.Spec); err != nil {
			return err
		}
		defer task.kubeCommon.SidecarVpnDelete(namespace, jobName)
	}

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

	podInfo, err := task.client.JobDescribe(namespace, jobName)
	if err != nil {
		return err
	}
	task.eventBus.Publish(newPodNameKubeEvent(namespace, podInfo.PodName, podInfo.ContainerName))

	// stop loader
	task.eventBus.Publish(newContainerWaitKubeLoaderEvent())

	logFileName := opts.GenerateLogFileName(taskModel.Kubernetes, podInfo.ContainerName)
	logOpts := &kubernetes.PodLogsOpts{
		Namespace: namespace,
		PodName:   podInfo.PodName,
		PodId:     podInfo.ContainerName,
	}
	task.eventBus.Publish(newPodLogKubeEvent(logFileName))
	if err := task.client.PodLogsTee(logOpts, logFileName); err != nil {
		return err
	}

	task.eventBus.Publish(newPodLogKubeConsoleEvent(logFileName))
	task.eventBus.Publish(newJobDeleteKubeEvent(namespace, jobName))
	return task.client.JobDelete(namespace, jobName)
}
