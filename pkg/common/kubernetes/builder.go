package kubernetes

import (
	"fmt"
	"path"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/util"
)

const (
	secretBasePath         = "/secrets"
	sidecarVpnTunnelVolume = "tun-device-volume"
	sidecarVpnTunnelPath   = "/dev/net/tun"
	sidecarVpnSecretVolume = "sidecar-vpn-volume"
	sidecarVpnSecretPath   = "openvpn/client.ovpn"
	sidecarVpnSecretKey    = "openvpn-config"
)

func buildSidecarVpnSecretName(containerName string) string {
	return fmt.Sprintf("%s-sidecar-vpn-secret", containerName)
}

func buildSidecarVpnSecret(namespace, containerName, secretValue string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      buildSidecarVpnSecretName(containerName),
			Namespace: namespace,
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{sidecarVpnSecretKey: []byte(secretValue)},
	}
}

func buildSidecarVpnContainer() corev1.Container {
	return corev1.Container{
		Name:            fmt.Sprintf("%svpn", commonModel.SidecarPrefixName),
		Image:           commonModel.SidecarVpnImageName,
		ImagePullPolicy: corev1.PullIfNotPresent,
		Env: []corev1.EnvVar{
			{Name: "OPENVPN_CONFIG", Value: path.Join(secretBasePath, sidecarVpnSecretPath)},
		},
		SecurityContext: &corev1.SecurityContext{
			Capabilities: &corev1.Capabilities{
				Add: []corev1.Capability{"NET_ADMIN"},
			},
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      sidecarVpnTunnelVolume,
				MountPath: sidecarVpnTunnelPath,
				ReadOnly:  true,
			},
			{
				Name:      sidecarVpnSecretVolume,
				MountPath: secretBasePath,
				ReadOnly:  true,
			},
		},
	}
}

func buildSidecarVpnVolumes(containerName string) []corev1.Volume {
	return []corev1.Volume{
		{
			Name: sidecarVpnTunnelVolume,
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: sidecarVpnTunnelPath,
				},
			},
		},
		{
			Name: sidecarVpnSecretVolume,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: buildSidecarVpnSecretName(containerName),
					Items: []corev1.KeyToPath{
						{Key: sidecarVpnSecretKey, Path: sidecarVpnSecretPath},
					},
				},
			},
		},
	}
}

func boolPtr(b bool) *bool { return &b }

func injectSidecarVpn(podSpec *corev1.PodSpec, mainContainerName string) {

	// https://kubernetes.io/docs/tasks/configure-pod-container/share-process-namespace
	//podSpec.ShareProcessNamespace = boolPtr(true)

	// disable ipv6, see https://kubernetes.io/docs/tasks/administer-cluster/sysctl-cluster
	podSpec.SecurityContext = &corev1.PodSecurityContext{
		Sysctls: []corev1.Sysctl{
			{Name: "net.ipv6.conf.all.disable_ipv6", Value: "0"},
		},
	}

	// inject containers
	podSpec.Containers = append(
		// order matters
		[]corev1.Container{
			buildSidecarVpnContainer(),
			// add fake sleep to allow sidecar-vpn to connect properly before starting the main container
			{
				Name:  fmt.Sprintf("%ssleep", commonModel.SidecarPrefixName),
				Image: "busybox",
				Lifecycle: &corev1.Lifecycle{
					PostStart: &corev1.LifecycleHandler{
						Exec: &corev1.ExecAction{
							Command: []string{"sleep", "1s"},
						},
					},
				},
			},
		},
		podSpec.Containers..., // current containers
	)

	// inject volumes
	podSpec.Volumes = append(
		podSpec.Volumes, // current volumes
		buildSidecarVpnVolumes(util.ToLowerKebabCase(mainContainerName))..., // join slices
	)
}
