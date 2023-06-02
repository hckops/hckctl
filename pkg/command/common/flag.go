package common

import (
	"github.com/thediveo/enumflag/v2"
)

const (
	NoneFlagShortHand = ""
)

type ProviderFlag enumflag.Flag

const (
	DockerFlag ProviderFlag = iota
	KubernetesFlag
	ArgoFlag
	CloudFlag
)

var ProviderIds = map[ProviderFlag][]string{
	DockerFlag:     {string(Docker)},
	KubernetesFlag: {string(Kubernetes)},
	ArgoFlag:       {string(Argo)},
	CloudFlag:      {string(Cloud)},
}
