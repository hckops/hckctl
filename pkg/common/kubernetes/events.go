package kubernetes

import (
	"fmt"

	"github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/event"
)

type kubeCommonEvent struct {
	kind  event.EventKind
	value string
}

func (e *kubeCommonEvent) Source() string {
	return model.KubernetesProvider
}

func (e *kubeCommonEvent) Kind() event.EventKind {
	return e.kind
}

func (e *kubeCommonEvent) String() string {
	return e.value
}

func newInitKubeClientEvent() *kubeCommonEvent {
	return &kubeCommonEvent{kind: event.LogDebug, value: "init kube client"}
}

func newCloseKubeClientEvent() *kubeCommonEvent {
	return &kubeCommonEvent{kind: event.LogDebug, value: "close kube client"}
}

func newSecretCreateKubeEvent(namespace string, name string) *kubeCommonEvent {
	return &kubeCommonEvent{kind: event.LogInfo, value: fmt.Sprintf("secret create: namespace=%s name=%s", namespace, name)}
}

func newSecretDeleteKubeEvent(namespace string, name string) *kubeCommonEvent {
	return &kubeCommonEvent{kind: event.LogInfo, value: fmt.Sprintf("secret delete: namespace=%s name=%s", namespace, name)}
}

func newSidecarVpnConnectKubeEvent(vpnName string) *kubeCommonEvent {
	return &kubeCommonEvent{kind: event.LogInfo, value: fmt.Sprintf("sidecar-vpn connect: vpnName=%s", vpnName)}
}

func newSidecarShareMountKubeEvent(shareDir string) *kubeCommonEvent {
	return &kubeCommonEvent{kind: event.LogInfo, value: fmt.Sprintf("sidecar-share mount: shareDir=%s", shareDir)}
}

func newSidecarShareUploadKubeEvent(localPath string, remotePath string) *kubeCommonEvent {
	return &kubeCommonEvent{kind: event.LogInfo, value: fmt.Sprintf("sidecar-share upload: localPath=%s remotePath=%s", localPath, remotePath)}
}

func newSidecarShareUploadKubeLoaderEvent() *kubeCommonEvent {
	return &kubeCommonEvent{kind: event.LoaderUpdate, value: fmt.Sprintf("uploading shared folder")}
}
