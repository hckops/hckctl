package model

type BoxProvider string

const (
	Docker     BoxProvider = "docker"
	Kubernetes BoxProvider = "kube"
	Cloud      BoxProvider = "cloud"
)

func (p BoxProvider) String() string {
	return string(p)
}
