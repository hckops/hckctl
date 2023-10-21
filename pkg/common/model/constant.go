package model

const (
	DockerProvider     = "docker"
	KubernetesProvider = "kube"
	CloudProvider      = "cloud"

	SidecarPrefixName             = "sidecar-"
	SidecarVpnImageName           = "hckops/alpine-openvpn:latest"
	SidecarVpnPrivilegedImageName = "hckops/alpine-openvpn-privileged:latest"
	SidecarShareImageName         = "busybox"
	SidecarShareDir               = "/hck/share"
)
