package kubernetes

import (
	"fmt"
	"path/filepath"
	"strings"

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
	sidecarShareVolume     = "sidecar-share-volume"
)

func buildSidecarVpnSecretName(podName string) string {
	return fmt.Sprintf("%s-sidecar-vpn-secret", util.ToLowerKebabCase(podName))
}

func buildSidecarVpnSecret(namespace, podName, secretValue string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      buildSidecarVpnSecretName(podName),
			Namespace: namespace,
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{sidecarVpnSecretKey: []byte(secretValue)},
	}
}

func buildSidecarVpnContainer(privileged bool) corev1.Container {

	// default
	imageName := commonModel.SidecarVpnImageName
	volumeMounts := []corev1.VolumeMount{
		{
			Name:      sidecarVpnSecretVolume,
			MountPath: secretBasePath,
			ReadOnly:  true,
		},
	}
	if privileged {
		imageName = commonModel.SidecarVpnPrivilegedImageName
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      sidecarVpnTunnelVolume,
			MountPath: sidecarVpnTunnelPath,
			ReadOnly:  true,
		})
	}

	return corev1.Container{
		Name:            fmt.Sprintf("%svpn", commonModel.SidecarPrefixName),
		Image:           imageName,
		ImagePullPolicy: corev1.PullIfNotPresent,
		Env: []corev1.EnvVar{
			{Name: "OPENVPN_CONFIG", Value: filepath.Join(secretBasePath, sidecarVpnSecretPath)},
		},
		SecurityContext: &corev1.SecurityContext{
			Capabilities: &corev1.Capabilities{
				Add: []corev1.Capability{"NET_ADMIN"},
			},
		},
		VolumeMounts: volumeMounts,
	}
}

func buildSidecarVpnVolumes(podName string, privileged bool) []corev1.Volume {
	volumes := []corev1.Volume{
		{
			Name: sidecarVpnSecretVolume,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: buildSidecarVpnSecretName(podName),
					Items: []corev1.KeyToPath{
						{Key: sidecarVpnSecretKey, Path: sidecarVpnSecretPath},
					},
				},
			},
		},
	}
	if privileged {
		volumes = append(volumes, corev1.Volume{
			Name: sidecarVpnTunnelVolume,
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: sidecarVpnTunnelPath,
				},
			},
		})
	}
	return volumes
}

func boolPtr(b bool) *bool { return &b }

func injectSidecarVpn(podSpec *corev1.PodSpec, podName string, privileged bool) {

	// https://kubernetes.io/docs/tasks/configure-pod-container/share-process-namespace
	//podSpec.ShareProcessNamespace = boolPtr(true)

	if privileged {
		// disable ipv6, see https://kubernetes.io/docs/tasks/administer-cluster/sysctl-cluster
		podSpec.SecurityContext = &corev1.PodSecurityContext{
			Sysctls: []corev1.Sysctl{
				{Name: "net.ipv6.conf.all.disable_ipv6", Value: "0"},
			},
		}
	}

	// inject containers
	podSpec.Containers = append(
		// order matters
		[]corev1.Container{
			buildSidecarVpnContainer(privileged),
			// add fake sleep to allow sidecar-vpn to connect properly before starting the main container
			{
				Name:  fmt.Sprintf("%ssleep", commonModel.SidecarPrefixName),
				Image: commonModel.SidecarShareImageName,
				Stdin: true, // fixes PostStartHookError
				Lifecycle: &corev1.Lifecycle{
					PostStart: &corev1.LifecycleHandler{
						Exec: &corev1.ExecAction{
							Command: []string{"sleep", "3s"},
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
		buildSidecarVpnVolumes(podName, privileged)..., // join slices
	)
}

func buildSidecarShareContainerName() string {
	return fmt.Sprintf("%sshare", commonModel.SidecarPrefixName)
}

func buildSidecarShareLock(remoteDir string) string {
	return filepath.Join(remoteDir, ".wait")
}

func buildSidecarShareContainer(remoteDir string) corev1.Container {
	return corev1.Container{
		Name:  buildSidecarShareContainerName(),
		Image: commonModel.SidecarShareImageName, // only requirement is the "tar" binary
		Stdin: true,
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      sidecarShareVolume,
				MountPath: remoteDir,
			},
		},
	}
}

func injectSidecarShare(podSpec *corev1.PodSpec, mainContainerName string, shareDir *commonModel.ShareDirInfo) {

	lockfile := buildSidecarShareLock(shareDir.RemotePath)

	if shareDir.LockDir {
		// create lockfile
		podSpec.InitContainers = append(
			podSpec.InitContainers,
			corev1.Container{
				Name:    "init-share-lock",
				Image:   commonModel.SidecarShareImageName,
				Command: []string{"touch", lockfile},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      sidecarShareVolume,
						MountPath: shareDir.RemotePath,
					},
				},
			},
		)
	}

	for index, c := range podSpec.Containers {
		if c.Name == mainContainerName {

			// ensure entrypoint is always empty: use only arguments
			podSpec.Containers[index].Command = []string{}

			// this is a hack, but at this time the only hooks available are postStart and preStop
			// block to allow the job to finish uploading before start the main container
			if shareDir.LockDir {
				arguments := podSpec.Containers[index].Args

				// override entrypoint and assume "sh" is always available
				podSpec.Containers[index].Command = []string{"/bin/sh", "-c"}

				// wait to start until lockfile is removed (upload is finished)
				preStartHook := fmt.Sprintf("while [ -f %s ]; do sleep 1; done", lockfile)
				podSpec.Containers[index].Args = []string{
					fmt.Sprintf("%s && %s", preStartHook, strings.Join(arguments, " ")),
				}
			}

			// mount read-only shared volume to main container
			podSpec.Containers[index].VolumeMounts = append(
				c.VolumeMounts,
				corev1.VolumeMount{
					Name:      sidecarShareVolume,
					MountPath: shareDir.RemotePath,
					ReadOnly:  true,
				},
			)
		}
	}

	// inject sidecar
	podSpec.Containers = append(
		podSpec.Containers, // current containers
		buildSidecarShareContainer(shareDir.RemotePath),
	)

	// inject shared volume between main container and sidecar
	podSpec.Volumes = append(
		podSpec.Volumes, // current volumes
		corev1.Volume{
			Name: sidecarShareVolume,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
	)
}
