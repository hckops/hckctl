package model

type TaskProvider string

const (
	Docker     TaskProvider = "docker"
	Kubernetes TaskProvider = "kube"
	Cloud      TaskProvider = "cloud"
)

func (p TaskProvider) String() string {
	return string(p)
}
