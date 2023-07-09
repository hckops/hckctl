package kubernetes

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/kubernetes"
	"github.com/hckops/hckctl/pkg/event"
)

type KubeBox struct {
	client       *kubernetes.KubeClient
	clientConfig *kubernetes.KubeClientConfig
	eventBus     *event.EventBus
}

func NewKubeBox(commonOpts *model.BoxCommonOptions, clientConfig *kubernetes.KubeClientConfig) (*KubeBox, error) {
	return newKubeBox(commonOpts, clientConfig)
}

func (box *KubeBox) Provider() model.BoxProvider {
	return model.Kubernetes
}

func (box *KubeBox) Events() *event.EventBus {
	return box.eventBus
}

func (box *KubeBox) Create(template *model.BoxV1) (*model.BoxInfo, error) {
	defer box.close()
	return box.createBox(template)
}

func (box *KubeBox) Connect(template *model.BoxV1, tunnelOpts *model.TunnelOptions, name string) error {
	defer box.client.Close()
	return box.connectBox(template, tunnelOpts, name)
}

func (box *KubeBox) Open(template *model.BoxV1, tunnelOpts *model.TunnelOptions) error {
	defer box.client.Close()
	return box.openBox(template, tunnelOpts)
}

func (box *KubeBox) Copy(string, string, string) error {
	defer box.client.Close()
	return errors.New("not implemented")
}

func (box *KubeBox) List() ([]model.BoxInfo, error) {
	defer box.client.Close()
	return box.listBoxes()
}

func (box *KubeBox) Delete(name string) error {
	defer box.client.Close()
	return box.deleteBox(name)
}

func (box *KubeBox) DeleteAll() ([]model.BoxInfo, error) {
	defer box.client.Close()
	return box.deleteBoxes()
}
