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

func newInitKubeClientEvent() *kubeTaskEvent {
	return &kubeTaskEvent{kind: event.LogDebug, value: "init kube client"}
}

func newCloseKubeClientEvent() *kubeTaskEvent {
	return &kubeTaskEvent{kind: event.LogDebug, value: "close kube client"}
}

func newNamespaceApplyKubeEvent(namespace string) *kubeTaskEvent {
	return &kubeTaskEvent{kind: event.LogInfo, value: fmt.Sprintf("namespace apply: namespace=%s", namespace)}
}
