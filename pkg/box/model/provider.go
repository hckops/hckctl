package model

import (
	"github.com/hckops/hckctl/pkg/common/model"
)

type BoxProvider string

const (
	Docker     BoxProvider = model.DockerProvider
	Kubernetes BoxProvider = model.KubernetesProvider
	Cloud      BoxProvider = model.CloudProvider
)

func (p BoxProvider) String() string {
	return string(p)
}
