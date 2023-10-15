package model

const (
	DockerProvider     = "docker"
	KubernetesProvider = "kube"
	CloudProvider      = "cloud"

	SidecarPrefixName   = "sidecar-"
	SidecarVpnImageName = "hckops/alpine-openvpn:latest"
	SidecarShareDir     = "/hck/share"
)
