package model

type BoxProvider uint

const (
	Docker BoxProvider = iota
	Kubernetes
	Cloud
)

var providerValue = []string{"docker", "kube", "cloud"}

func (p BoxProvider) String() string {
	return providerValue[p]
}
