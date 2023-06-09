package kubernetes

import (
	"github.com/hckops/hckctl/pkg/box"
	"github.com/hckops/hckctl/pkg/template/model"
)

type KubeClientOpts struct {
	template model.BoxV1
}

type KubeClient struct {
	opts *KubeClientOpts
}

func (client *KubeClient) Open() (*box.Connection, error) {
	return nil, nil
}
