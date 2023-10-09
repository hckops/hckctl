package kubernetes

import (
	"fmt"

	"github.com/hckops/hckctl/pkg/event"
	"github.com/hckops/hckctl/pkg/task/model"
)

type kubeTaskEvent struct {
	kind  event.EventKind
	value string
}

func (e *kubeTaskEvent) Source() string {
	return model.Kubernetes.String()
}

func (e *kubeTaskEvent) Kind() event.EventKind {
	return e.kind
}

func (e *kubeTaskEvent) String() string {
	return e.value
}

func newNamespaceApplyKubeEvent(namespace string) *kubeTaskEvent {
	return &kubeTaskEvent{kind: event.LogInfo, value: fmt.Sprintf("namespace apply: namespace=%s", namespace)}
}

func newJobCreateStatusKubeEvent(status string) *kubeTaskEvent {
	return &kubeTaskEvent{kind: event.LogDebug, value: status}
}

func newJobCreateKubeEvent(namespace string, name string) *kubeTaskEvent {
	return &kubeTaskEvent{kind: event.LogInfo, value: fmt.Sprintf("job create: namespace=%s name=%s", namespace, name)}
}

func newJobDeleteKubeEvent(namespace string, name string) *kubeTaskEvent {
	return &kubeTaskEvent{kind: event.LogInfo, value: fmt.Sprintf("job delete: namespace=%s name=%s", namespace, name)}
}

func newPodNameKubeEvent(namespace string, podName string, containerName string) *kubeTaskEvent {
	return &kubeTaskEvent{kind: event.LogInfo, value: fmt.Sprintf("found unique pod: namespace=%s podName=%s containerName=%s", namespace, podName, containerName)}
}

func newPodLogKubeEvent(logFileName string) *kubeTaskEvent {
	return &kubeTaskEvent{kind: event.LogInfo, value: fmt.Sprintf("pod log: logFileName=%s", logFileName)}
}

func newPodLogKubeConsoleEvent(logFileName string) *kubeTaskEvent {
	return &kubeTaskEvent{kind: event.PrintConsole, value: fmt.Sprintf("output file: %s", logFileName)}
}

func newContainerWaitKubeLoaderEvent() *kubeTaskEvent {
	return &kubeTaskEvent{kind: event.LoaderStop, value: "waiting"}
}
