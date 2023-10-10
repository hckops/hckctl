package kubernetes

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	commonModel "github.com/hckops/hckctl/pkg/common/model"
)

const (
	secretBasePath         = "/secrets"
	sidecarVpnTunnelVolume = "tun-device-volume"
	sidecarVpnTunnelPath   = "/dev/net/tun"
	sidecarVpnSecretVolume = "sidecar-vpn-volume"
	sidecarVpnSecretPath   = "/secrets/openvpn/client.ovpn"
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
		Name:            "sidecar-vpn",
		Image:           commonModel.SidecarVpnImageName,
		ImagePullPolicy: corev1.PullIfNotPresent,
		Env: []corev1.EnvVar{
			{Name: "OPENVPN_CONFIG", Value: sidecarVpnSecretPath},
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
