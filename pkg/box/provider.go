package box

type BoxProvider string

const (
	Docker     BoxProvider = "docker"
	Kubernetes BoxProvider = "kube"
	Argo       BoxProvider = "argo"
	Cloud      BoxProvider = "cloud"
)

func BoxProviderValues() []string {
	return []string{string(Docker), string(Kubernetes), string(Argo), string(Cloud)}
}
