package box

import (
	"github.com/hckops/hckctl/pkg/template/model"
)

type KubeClientOpts struct {
	template model.BoxV1
}

type KubeClient struct {
	opts *KubeClientOpts
}

func (client *KubeClient) Open() (*Connection, error) {
	return nil, nil
}
