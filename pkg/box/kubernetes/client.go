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
	streams      *model.BoxStreams
	eventBus     *event.EventBus
}

func NewKubeBox(internalOpts *model.BoxInternalOpts, kubeConfig *kubernetes.KubeClientConfig) (*KubeBox, error) {
	return newKubeBox(internalOpts, kubeConfig)
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

func (box *KubeBox) Exec(template *model.BoxV1, name string) error {
	defer box.client.Close()
	return box.execBox(template, name)
}

func (box *KubeBox) Tunnel(string) error {
	defer box.client.Close()
	return errors.New("not implemented")
}

func (box *KubeBox) Open(template *model.BoxV1) error {
	defer box.client.Close()
	return box.openBox(template)
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
