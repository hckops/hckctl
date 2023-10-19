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

func newDeploymentListKubeEvent(index int, namespace string, name string, healthy bool) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LogDebug, value: fmt.Sprintf("deployment list: (%d) namespace=%s name=%s healthy=%v", index, namespace, name, healthy)}
}

func newDeploymentDescribeKubeEvent(namespace string, name string) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LogInfo, value: fmt.Sprintf("deployment describe: namespace=%s name=%s", namespace, name)}
}

func newDeploymentDeleteKubeEvent(namespace string, name string) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LogInfo, value: fmt.Sprintf("deployment delete: namespace=%s name=%s", namespace, name)}
}

func newPodNameKubeEvent(namespace string, name string, containerName string) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LogInfo, value: fmt.Sprintf("found unique pod: namespace=%s name=%s containerName=%s", namespace, name, containerName)}
}

func newPodExecKubeEvent(templateName string, namespace string, name string, command string) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LogInfo, value: fmt.Sprintf("pod attach: templateName=%s namespace=%s name=%s command=%s", templateName, namespace, name, command)}
}

func newPodExecKubeLoaderEvent() *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LoaderStop, value: "waiting"}
}

func newPodLogsKubeEvent(namespace string, name string) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LogDebug, value: fmt.Sprintf("pod logs: namespace=%s name=%s", namespace, name)}
}

func newPodLogsExitKubeEvent(namespace string, name string) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LogDebug, value: fmt.Sprintf("pod logs exit: namespace=%s name=%s", namespace, name)}
}

func newPodLogsExitKubeConsoleEvent() *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LoaderUpdate, value: "killing"}
}

func newPodLogsErrorKubeEvent(namespace string, name string, err error) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LogError, value: fmt.Sprintf("pod logs error: namespace=%s name=%s error=%v", namespace, name, err)}
}

func newPodPortForwardIgnoreKubeEvent(namespace string, name string) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LogWarning, value: fmt.Sprintf("pod port-forward ignored: namespace=%s name=%s", namespace, name)}
}

func newPodPortForwardBindingKubeEvent(namespace, name string, port model.BoxPort) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LogInfo, value: fmt.Sprintf(
		"pod port-forward: namespace=%s name=%s portAlias=%s portRemote=%s portLocal=%s",
		namespace, name, port.Alias, port.Remote, port.Local)}
}

func newPodPortForwardBindingKubeConsoleEvent(namespace string, containerName string, port model.BoxPort, padding int) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.PrintConsole, value: fmt.Sprintf(
		"[%s/%s][%-*s] tunnel (remote) %s -> (local) %s",
		namespace, containerName, padding, port.Alias, port.Remote, port.Local)}
}

func newPodPortForwardErrorKubeEvent(namespace string, name string, err error) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LogError, value: fmt.Sprintf("pod port-forward error: namespace=%s name=%s error=%v", namespace, name, err)}
}

func newPodEnvKubeEvent(namespace string, name string, env model.BoxEnv) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.LogInfo, value: fmt.Sprintf("pod env: namespace=%s name=%s key=%s value=%s", namespace, name, env.Key, env.Value)}
}

func newPodEnvKubeConsoleEvent(namespace string, containerName string, env model.BoxEnv) *kubeBoxEvent {
	return &kubeBoxEvent{kind: event.PrintConsole, value: fmt.Sprintf("[%s/%s] %s=%s", namespace, containerName, env.Key, env.Value)}
}
