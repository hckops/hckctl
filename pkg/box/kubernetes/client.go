package kubernetes

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/kubernetes"
	"github.com/hckops/hckctl/pkg/event"
	"github.com/hckops/hckctl/pkg/provider"
)

type KubeBoxClient struct {
	client     *kubernetes.KubeClient
	clientOpts *provider.KubeOptions
	eventBus   *event.EventBus
}

func NewKubeBoxClient(commonOpts *model.CommonBoxOptions, kubeOpts *provider.KubeOptions) (*KubeBoxClient, error) {
	return newKubeBoxClient(commonOpts, kubeOpts)
}

func (box *KubeBoxClient) Provider() model.BoxProvider {
	return model.Kubernetes
}

func (box *KubeBoxClient) Events() *event.EventBus {
	return box.eventBus
}

func (box *KubeBoxClient) Create(opts *model.CreateOptions) (*model.BoxInfo, error) {
	defer box.close()
	return box.createBox(opts)
}

func (box *KubeBoxClient) Connect(opts *model.ConnectOptions) error {
	defer box.close()
	return box.connectBox(opts)
}

func (box *KubeBoxClient) Copy(string, string, string) error {
	defer box.close()
	return errors.New("not implemented")
}

func (box *KubeBoxClient) Describe(name string) (*model.BoxDetails, error) {
	defer box.close()
	return box.describeBox(name)
}

func (box *KubeBoxClient) List() ([]model.BoxInfo, error) {
	defer box.close()
	return box.listBoxes()
}

func (box *KubeBoxClient) Delete(names []string) ([]string, error) {
	defer box.close()
	return box.deleteBoxes(names)
}

func (box *KubeBoxClient) Clean() error {
	defer box.close()
	return box.clean()
}

func (box *KubeBoxClient) Version() (string, error) {
	return "", errors.New("not implemented")
}
