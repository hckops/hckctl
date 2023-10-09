package kubernetes

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/client/kubernetes"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/event"
)

type KubeCommonClient struct {
	client     *kubernetes.KubeClient
	clientOpts *commonModel.KubeOptions
	eventBus   *event.EventBus
}

func NewKubeCommonClient(kubeOpts *commonModel.KubeOptions, eventBus *event.EventBus) (*KubeCommonClient, error) {
	eventBus.Publish(newInitKubeClientEvent())

	kubeClient, err := kubernetes.NewKubeClient(kubeOpts.InCluster, kubeOpts.ConfigPath)
	if err != nil {
		return nil, errors.Wrap(err, "error kube common client")
	}

	return &KubeCommonClient{
		client:     kubeClient,
		clientOpts: kubeOpts,
		eventBus:   eventBus,
	}, nil
}

func (common *KubeCommonClient) GetClient() *kubernetes.KubeClient {
	return common.client
}

func (common *KubeCommonClient) Close() error {
	common.eventBus.Publish(newCloseKubeClientEvent())
	common.eventBus.Close()
	return common.client.Close()
}
