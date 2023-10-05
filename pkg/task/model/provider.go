package model

import (
	"github.com/hckops/hckctl/pkg/common/model"
)

type TaskProvider string

const (
	Docker     TaskProvider = model.DockerProvider
	Kubernetes TaskProvider = model.KubernetesProvider
	Cloud      TaskProvider = model.CloudProvider
)

func (p TaskProvider) String() string {
	return string(p)
}
