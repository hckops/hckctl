package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/hckops/hckctl/pkg/client/kubernetes"
	"github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/util"
)

func newPodSpecTest(containerName string) *corev1.Pod {
	return &corev1.Pod{
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    containerName,
					Image:   "my-image",
					Command: []string{"xyz", "abc"},
					Args:    []string{"foo", "bar"},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: "my-volume",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "my-path",
						},
					},
				},
			},
		},
	}
}

func TestBuildSidecarVpnSecret(t *testing.T) {

	expected := `
apiVersion: v1
data:
  openvpn-config: bXktdmFsdWU=
kind: Secret
metadata:
  creationTimestamp: null
  name: my-container-name-sidecar-vpn-secret
  namespace: my-namespace
type: Opaque
`

	actual := buildSidecarVpnSecret("my-namespace", "my-container-name", "my-value")
	// fix model
	actual.TypeMeta = metav1.TypeMeta{Kind: "Secret", APIVersion: "v1"}

	assert.YAMLEqf(t, expected, kubernetes.ObjectToYaml(actual), "unexpected secret")

	decoded, ok := util.Base64Decode("bXktdmFsdWU=")
	assert.True(t, ok)
	assert.Equal(t, "my-value", decoded)
}

func TestBuildSidecarVpnPod(t *testing.T) {

	expected := `
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
spec:
  containers:
  - env:
    - name: OPENVPN_CONFIG
      value: /secrets/openvpn/client.ovpn
    image: hckops/alpine-openvpn:latest
    imagePullPolicy: IfNotPresent
    name: sidecar-vpn
    resources: {}
    securityContext:
      capabilities:
        add:
        - NET_ADMIN
    volumeMounts:
    - mountPath: /dev/net/tun
      name: tun-device-volume
      readOnly: true
    - mountPath: /secrets
      name: sidecar-vpn-volume
      readOnly: true
  volumes:
  - hostPath:
      path: /dev/net/tun
    name: tun-device-volume
  - name: sidecar-vpn-volume
    secret:
      items:
      - key: openvpn-config
        path: openvpn/client.ovpn
      secretName: main-container-sidecar-vpn-secret
status: {}
`

	actual := &corev1.Pod{
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{buildSidecarVpnContainer()},
			Volumes:    buildSidecarVpnVolumes("main-container"),
		},
	}
	// fix model
	actual.TypeMeta = metav1.TypeMeta{Kind: "Pod", APIVersion: "v1"}

	assert.YAMLEqf(t, expected, kubernetes.ObjectToYaml(actual), "unexpected pod")
}

func TestInjectSidecarVpn(t *testing.T) {

	expected := `
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
spec:
  containers:
  - env:
    - name: OPENVPN_CONFIG
      value: /secrets/openvpn/client.ovpn
    image: hckops/alpine-openvpn:latest
    imagePullPolicy: IfNotPresent
    name: sidecar-vpn
    resources: {}
    securityContext:
      capabilities:
        add:
        - NET_ADMIN
    volumeMounts:
    - mountPath: /dev/net/tun
      name: tun-device-volume
      readOnly: true
    - mountPath: /secrets
      name: sidecar-vpn-volume
      readOnly: true
  - image: busybox
    lifecycle:
      postStart:
        exec:
          command:
          - sleep
          - 1s
    name: sidecar-sleep
    resources: {}
    stdin: true
  - args:
    - foo
    - bar
    command:
    - xyz
    - abc
    image: my-image
    name: my-name
    resources: {}
  securityContext:
    sysctls:
    - name: net.ipv6.conf.all.disable_ipv6
      value: "0"
  volumes:
  - hostPath:
      path: my-path
    name: my-volume
  - hostPath:
      path: /dev/net/tun
    name: tun-device-volume
  - name: sidecar-vpn-volume
    secret:
      items:
      - key: openvpn-config
        path: openvpn/client.ovpn
      secretName: my-name-sidecar-vpn-secret
status: {}
`

	containerName := "my-name"
	actual := newPodSpecTest(containerName)
	injectSidecarVpn(&actual.Spec, containerName)
	// fix model
	actual.TypeMeta = metav1.TypeMeta{Kind: "Pod", APIVersion: "v1"}

	assert.YAMLEqf(t, expected, kubernetes.ObjectToYaml(actual), "unexpected pod")
}

func TestInjectSidecarShare(t *testing.T) {

	expected := `
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
spec:
  containers:
  - args:
    - foo
    - bar
    image: my-image
    name: my-name
    resources: {}
    volumeMounts:
    - mountPath: /tmp/foo
      name: sidecar-share-volume
      readOnly: true
  - image: busybox
    name: sidecar-share
    resources: {}
    stdin: true
    volumeMounts:
    - mountPath: /tmp/foo
      name: sidecar-share-volume
  volumes:
  - hostPath:
      path: my-path
    name: my-volume
  - emptyDir: {}
    name: sidecar-share-volume
status: {}
`

	containerName := "my-name"
	shareDir := &model.ShareDirInfo{RemotePath: "/tmp/foo", LockDir: false}
	actual := newPodSpecTest(containerName)
	injectSidecarShare(&actual.Spec, containerName, shareDir)
	// fix model
	actual.TypeMeta = metav1.TypeMeta{Kind: "Pod", APIVersion: "v1"}

	assert.YAMLEqf(t, expected, kubernetes.ObjectToYaml(actual), "unexpected pod")
}

func TestInjectSidecarShareLockDir(t *testing.T) {

	expected := `
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
spec:
  containers:
  - args:
    - while [ -f /tmp/foo/.wait ]; do sleep 1; done && foo bar
    command:
    - /bin/sh
    - -c
    image: my-image
    name: my-name
    resources: {}
    volumeMounts:
    - mountPath: /tmp/foo
      name: sidecar-share-volume
      readOnly: true
  - image: busybox
    name: sidecar-share
    resources: {}
    stdin: true
    volumeMounts:
    - mountPath: /tmp/foo
      name: sidecar-share-volume
  initContainers:
  - command:
    - touch
    - /tmp/foo/.wait
    image: busybox
    name: init-share-lock
    resources: {}
    volumeMounts:
    - mountPath: /tmp/foo
      name: sidecar-share-volume
  volumes:
  - hostPath:
      path: my-path
    name: my-volume
  - emptyDir: {}
    name: sidecar-share-volume
status: {}
`

	containerName := "my-name"
	shareDir := &model.ShareDirInfo{RemotePath: "/tmp/foo", LockDir: true}
	actual := newPodSpecTest(containerName)
	injectSidecarShare(&actual.Spec, containerName, shareDir)
	// fix model
	actual.TypeMeta = metav1.TypeMeta{Kind: "Pod", APIVersion: "v1"}

	assert.YAMLEqf(t, expected, kubernetes.ObjectToYaml(actual), "unexpected pod")
}
