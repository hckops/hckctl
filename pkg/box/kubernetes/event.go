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

func newResourcesDeployKubeLoaderEvent(namespace string, name string) *kubeEvent {
	return &kubeEvent{kind: event.LoaderUpdate, value: fmt.Sprintf("deploying %s/%s", namespace, name)}
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

func newPodNameKubeEvent(namespace string, podName string, containerName string) *kubeEvent {
	return &kubeEvent{kind: event.LogDebug, value: fmt.Sprintf("found unique pod: namespace=%s podName=%s containerName=%s", namespace, podName, containerName)}
}

func newPodPortForwardSkippedKubeEvent(namespace string, podId string) *kubeEvent {
	return &kubeEvent{kind: event.LogDebug, value: fmt.Sprintf("pod port-forward skipped: namespace=%s podId=%s", namespace, podId)}
}

func newPodPortForwardBindingKubeEvent(namespace, podId string, port model.BoxPort) *kubeEvent {
	return &kubeEvent{kind: event.LogDebug, value: fmt.Sprintf(
		"pod port-forward: namespace=%s podId=%s portAlias=%s portRemote=%s portLocal=%s",
		namespace, podId, port.Alias, port.Remote, port.Local)}
}

func newPodPortForwardBindingKubeConsoleEvent(namespace string, podName string, port model.BoxPort, padding int) *kubeEvent {
	return &kubeEvent{kind: event.PrintConsole, value: fmt.Sprintf(
		"[%s/%s][%-*s] tunnel (remote) %s -> (local) %s",
		namespace, podName, padding, port.Alias, port.Remote, port.Local)}
}

func newPodPortForwardErrorKubeEvent(namespace string, podId string, err error) *kubeEvent {
	return &kubeEvent{kind: event.LogWarning, value: fmt.Sprintf("pod port-forward error: namespace=%s podId=%s error=%v", namespace, podId, err)}
}

func newPodExecKubeEvent(templateName string, namespace string, podId string, command string) *kubeEvent {
	return &kubeEvent{kind: event.LogDebug, value: fmt.Sprintf("pod attach: templateName=%s namespace=%s podId=%s command=%s", templateName, namespace, podId, command)}
}

func newPodExecKubeLoaderEvent() *kubeEvent {
	return &kubeEvent{kind: event.LoaderStop, value: "waiting"}
}

func newDeploymentListKubeEvent(index int, namespace string, deploymentName string, podId string, healthy bool) *kubeEvent {
	return &kubeEvent{kind: event.LogDebug, value: fmt.Sprintf("deployment list: (%d) namespace=%s deploymentName=%s podId=%s healthy=%v", index, namespace, deploymentName, podId, healthy)}
}

func newDeploymentSearchKubeEvent(namespace string, name string) *kubeEvent {
	return &kubeEvent{kind: event.LogDebug, value: fmt.Sprintf("deployment search: namespace=%s name=%s", namespace, name)}
}
