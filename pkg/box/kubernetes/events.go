package kubernetes

import (
	"fmt"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/event"
)

type kubeBoxEvent struct {
	kind  event.EventKind
	value string
}

func (e *kubeBoxEvent) Source() string {
	return model.Kubernetes.String()
}

func (e *kubeBoxEvent) Kind() event.EventKind {
	return e.kind
}

func (e *kubeBoxEvent) String() string {
	return e.value
}

func newNamespaceApplyKubeEvent(namespace string) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LogInfo, value: fmt.Sprintf("namespace apply: namespace=%s", namespace)}
}

func newNamespaceDeleteKubeEvent(namespace string) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LogInfo, value: fmt.Sprintf("namespace delete: namespace=%s", namespace)}
}

func newResourcesDeployKubeLoaderEvent(namespace string, name string) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LoaderUpdate, value: fmt.Sprintf("deploying %s/%s", namespace, name)}
}

func newResourcesDeleteIgnoreKubeEvent(namespace string, name string) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LogWarning, value: fmt.Sprintf("resources delete ignored: namespace=%s name=%s", namespace, name)}
}

func newServiceCreateKubeEvent(namespace string, name string) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LogInfo, value: fmt.Sprintf("service create: namespace=%s name=%s", namespace, name)}
}

func newServiceCreateIgnoreKubeEvent(namespace string, name string) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LogWarning, value: fmt.Sprintf("service create ignored: namespace=%s name=%s", namespace, name)}
}

func newServiceDescribeKubeEvent(namespace string, name string) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LogInfo, value: fmt.Sprintf("service describe: namespace=%s name=%s", namespace, name)}
}

func newServiceDeleteKubeEvent(namespace string, name string) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LogInfo, value: fmt.Sprintf("service delete: namespace=%s name=%s", namespace, name)}
}

func newDeploymentCreateKubeEvent(namespace string, name string) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LogInfo, value: fmt.Sprintf("deployment create: namespace=%s name=%s", namespace, name)}
}

func newDeploymentCreateStatusKubeEvent(status string) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LogDebug, value: status}
}

func newDeploymentSearchKubeEvent(namespace string, name string) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LogInfo, value: fmt.Sprintf("deployment search: namespace=%s name=%s", namespace, name)}
}

func newDeploymentListKubeEvent(index int, namespace string, deploymentName string, podId string, healthy bool) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LogDebug, value: fmt.Sprintf("deployment list: (%d) namespace=%s deploymentName=%s podId=%s healthy=%v", index, namespace, deploymentName, podId, healthy)}
}

func newDeploymentDescribeKubeEvent(namespace string, name string) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LogInfo, value: fmt.Sprintf("deployment describe: namespace=%s name=%s", namespace, name)}
}

func newDeploymentDeleteKubeEvent(namespace string, name string) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LogInfo, value: fmt.Sprintf("deployment delete: namespace=%s name=%s", namespace, name)}
}

func newPodNameKubeEvent(namespace string, podName string, containerName string) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LogInfo, value: fmt.Sprintf("found unique pod: namespace=%s podName=%s containerName=%s", namespace, podName, containerName)}
}

func newPodExecKubeEvent(templateName string, namespace string, podId string, command string) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LogInfo, value: fmt.Sprintf("pod attach: templateName=%s namespace=%s podId=%s command=%s", templateName, namespace, podId, command)}
}

func newPodExecKubeLoaderEvent() *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LoaderStop, value: "waiting"}
}

func newPodPortForwardIgnoreKubeEvent(namespace string, podId string) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LogWarning, value: fmt.Sprintf("pod port-forward ignored: namespace=%s podId=%s", namespace, podId)}
}

func newPodPortForwardBindingKubeEvent(namespace, podId string, port model.BoxPort) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LogInfo, value: fmt.Sprintf(
		"pod port-forward: namespace=%s podId=%s portAlias=%s portRemote=%s portLocal=%s",
		namespace, podId, port.Alias, port.Remote, port.Local)}
}

func newPodPortForwardBindingKubeConsoleEvent(namespace string, podName string, port model.BoxPort, padding int) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.PrintConsole, value: fmt.Sprintf(
		"[%s/%s][%-*s] tunnel (remote) %s -> (local) %s",
		namespace, podName, padding, port.Alias, port.Remote, port.Local)}
}

func newPodPortForwardErrorKubeEvent(namespace string, podId string, err error) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LogError, value: fmt.Sprintf("pod port-forward error: namespace=%s podId=%s error=%v", namespace, podId, err)}
}

func newPodEnvKubeEvent(namespace string, podId string, env model.BoxEnv) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LogInfo, value: fmt.Sprintf("pod env: namespace=%s podId=%s key=%s value=%s", namespace, podId, env.Key, env.Value)}
}

func newPodEnvKubeConsoleEvent(namespace string, podId string, env model.BoxEnv) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.PrintConsole, value: fmt.Sprintf("[%s/%s] %s=%s", namespace, podId, env.Key, env.Value)}
}
