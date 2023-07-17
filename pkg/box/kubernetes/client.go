package kubernetes

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/kubernetes"
	"github.com/hckops/hckctl/pkg/event"
)

type KubeBoxClient struct {
	client     *kubernetes.KubeClient
	clientOpts *model.KubeBoxOptions
	eventBus   *event.EventBus
}

func NewKubeBoxClient(commonOpts *model.CommonBoxOptions, kubeOpts *model.KubeBoxOptions) (*KubeBoxClient, error) {
	return newKubeBoxClient(commonOpts, kubeOpts)
}

func (box *KubeBoxClient) Provider() model.BoxProvider {
	return model.Kubernetes
}

func (box *KubeBoxClient) Events() *event.EventBus {
	return box.eventBus
}

func (box *KubeBoxClient) Create(templateOpts *model.TemplateOptions) (*model.BoxInfo, error) {
	defer box.close()
	return box.createBox(templateOpts)
}

func (box *KubeBoxClient) Connect(template *model.BoxV1, tunnelOpts *model.TunnelOptions, name string) error {
	defer box.close()
	return box.connectBox(template, tunnelOpts, name)
}

func (box *KubeBoxClient) Open(templateOpts *model.TemplateOptions, tunnelOpts *model.TunnelOptions) error {
	defer box.close()
	return box.openBox(templateOpts, tunnelOpts)
}

func (box *KubeBoxClient) Copy(string, string, string) error {
	defer box.close()
	return errors.New("not implemented")
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
