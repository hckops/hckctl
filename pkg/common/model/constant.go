package model

const (
	DockerProvider     = "docker"
	KubernetesProvider = "kube"
	CloudProvider      = "cloud"

	SidecarVpnImageName = "hckops/alpine-openvpn:latest"
	MountShareDir       = "/hck/share"
)
