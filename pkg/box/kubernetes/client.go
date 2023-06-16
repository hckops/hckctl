package kubernetes

import (
	"github.com/hckops/hckctl/pkg/client/kubernetes"
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/event"
)

type KubeBox struct {
	client *kubernetes.KubeClient
	opts   *model.BoxOpts
}

func NewKubeBox(opts *model.BoxOpts) (*KubeBox, error) {
	return newKubeBox(opts)
}

func (box *KubeBox) Events() *event.EventBus {
	return nil
}

func (box *KubeBox) Create(template *model.BoxV1) (*model.BoxInfo, error) {
	defer box.close()
	return nil, errors.New("not implemented")
}

func (box *KubeBox) Exec(name string, command string) error {
	defer box.client.Close()
	return errors.New("not implemented")
}

func (box *KubeBox) Open(template *model.BoxV1) error {
	defer box.client.Close()
	return errors.New("not implemented")
}

func (box *KubeBox) List() ([]model.BoxInfo, error) {
	defer box.client.Close()
	return nil, errors.New("not implemented")
}

func (box *KubeBox) Copy(string, string, string) error {
	defer box.client.Close()
	return errors.New("not implemented")
}

func (box *KubeBox) Tunnel(string) error {
	defer box.client.Close()
	return errors.New("not supported")
}

func (box *KubeBox) Delete(name string) error {
	defer box.client.Close()
	return errors.New("not implemented")
}

func (box *KubeBox) DeleteAll() ([]model.BoxInfo, error) {
	defer box.client.Close()
	return nil, errors.New("not implemented")
}
