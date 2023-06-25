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

func NewKubeBox(internalOpts *model.BoxInternalOptions, clientConfig *kubernetes.KubeClientConfig) (*KubeBox, error) {
	return newKubeBox(internalOpts, clientConfig)
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
	// TODO tunnelOpts
	return box.execBox(template, name)
}

func (box *KubeBox) Open(template *model.BoxV1, tunnelOpts *model.TunnelOptions) error {
	defer box.client.Close()
	// TODO tunnelOpts
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
