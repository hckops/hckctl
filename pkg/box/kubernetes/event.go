package kubernetes

import (
	"fmt"
	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/event"
)

type kubeEvent struct {
	kind  event.EventKind
	value string
}

func (e *kubeEvent) Source() string {
	return model.Kubernetes.String()
}

func (e *kubeEvent) Kind() event.EventKind {
	return e.kind
}

func (e *kubeEvent) String() string {
	return e.value
}

func newClientInitKubeEvent() *kubeEvent {
	return &kubeEvent{kind: event.LogDebug, value: "init kube client"}
}

func newClientCloseKubeEvent() *kubeEvent {
	return &kubeEvent{kind: event.LogDebug, value: "close kube client"}
}

func newNamespaceApplyKubeEvent(namespace string) *kubeEvent {
	return &kubeEvent{kind: event.LogDebug, value: fmt.Sprintf("namespace apply: namespace=%s", namespace)}
}

func newNamespaceDeleteKubeEvent(namespace string) *kubeEvent {
	return &kubeEvent{kind: event.LogDebug, value: fmt.Sprintf("namespace delete: namespace=%s", namespace)}
}

func newNamespaceDeleteSkippedKubeEvent(namespace string) *kubeEvent {
	return &kubeEvent{kind: event.LogWarning, value: fmt.Sprintf("namespace delete skipped: namespace=%s", namespace)}
}

func newResourcesCreateLoaderKubeEvent(namespace string, name string) *kubeEvent {
	return &kubeEvent{kind: event.LoaderUpdate, value: fmt.Sprintf("creating %s/%s", namespace, name)}
}

func newResourcesDeleteSkippedKubeEvent(namespace string, name string) *kubeEvent {
	return &kubeEvent{kind: event.LogWarning, value: fmt.Sprintf("resources delete skipped: namespace=%s name=%s", namespace, name)}
}

func newServiceCreateKubeEvent(namespace string, name string) *kubeEvent {
	return &kubeEvent{kind: event.LogDebug, value: fmt.Sprintf("service create: namespace=%s name=%s", namespace, name)}
}

func newServiceCreateSkippedKubeEvent(namespace string, name string) *kubeEvent {
	return &kubeEvent{kind: event.LogDebug, value: fmt.Sprintf("service create skipped: namespace=%s name=%s", namespace, name)}
}

func newServiceDeleteKubeEvent(namespace string, name string) *kubeEvent {
	return &kubeEvent{kind: event.LogDebug, value: fmt.Sprintf("service delete: namespace=%s name=%s", namespace, name)}
}

func newDeploymentCreateKubeEvent(namespace string, name string) *kubeEvent {
	return &kubeEvent{kind: event.LogDebug, value: fmt.Sprintf("deployment create: namespace=%s name=%s", namespace, name)}
}

func newDeploymentCreateStatusKubeEvent(status string) *kubeEvent {
	return &kubeEvent{kind: event.LogDebug, value: status}
}

func newDeploymentDeleteKubeEvent(namespace string, name string) *kubeEvent {
	return &kubeEvent{kind: event.LogDebug, value: fmt.Sprintf("deployment delete: namespace=%s name=%s", namespace, name)}
}

func newPodNameKubeEvent(namespace string, name string) *kubeEvent {
	return &kubeEvent{kind: event.LogDebug, value: fmt.Sprintf("found unique pod: namespace=%s name=%s", namespace, name)}
}

func newPodPortForwardSkippedKubeEvent(namespace string, name string) *kubeEvent {
	return &kubeEvent{kind: event.LogDebug, value: fmt.Sprintf("pod port-forward skipped: namespace=%s name=%s", namespace, name)}
}

func newPodPortForwardBindingKubeEvent(namespace, name string, port model.BoxPort) *kubeEvent {
	return &kubeEvent{kind: event.LogDebug, value: fmt.Sprintf(
		"pod port-forward: namespace=%s name=%s portAlias=%s portRemote=%s portLocal=%s",
		namespace, name, port.Alias, port.Remote, port.Local)}
}

func newPodPortForwardBindingConsoleKubeEvent(namespace string, name string, port model.BoxPort) *kubeEvent {
	return &kubeEvent{kind: event.PrintConsole, value: fmt.Sprintf(
		"[%s/%s][%s]   \texpose (remote) %s -> (local) %s",
		namespace, name, port.Alias, port.Remote, port.Local)}
}

func newPodPortForwardErrorKubeEvent(namespace string, name string, err error) *kubeEvent {
	return &kubeEvent{kind: event.LogWarning, value: fmt.Sprintf("pod port-forward error: namespace=%s name=%s error=%v", namespace, name, err)}
}

func newDeploymentListKubeEvent(index int, namespace string, deploymentName string, podName string) *kubeEvent {
	return &kubeEvent{kind: event.LogDebug, value: fmt.Sprintf("deployment list: (%d) namespace=%s  deploymentName=%s podName=%s", index, namespace, deploymentName, podName)}
}
