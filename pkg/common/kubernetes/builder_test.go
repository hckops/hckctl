package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/hckops/hckctl/pkg/client/kubernetes"
	"github.com/hckops/hckctl/pkg/util"
)

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
        path: /secrets/openvpn/client.ovpn
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
