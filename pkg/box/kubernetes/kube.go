package kubernetes

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/kubernetes"
)

func newKubeBox(opts *model.BoxOpts) (*KubeBox, error) {
	opts.EventBus.Publish(newClientInitKubeEvent())

	kubeClient, err := kubernetes.NewKubeClient()
	if err != nil {
		return nil, errors.Wrap(err, "error kube box")
	}

	return &KubeBox{
		client: kubeClient,
		opts:   opts,
	}, nil
}

func (box *KubeBox) close() error {
	box.opts.EventBus.Publish(newClientCloseKubeEvent())
	box.opts.EventBus.Close()
	return box.client.Close()
}
