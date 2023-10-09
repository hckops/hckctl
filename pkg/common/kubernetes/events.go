package kubernetes

import (
	"github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/event"
)

type kubeTaskEvent struct {
	kind  event.EventKind
	value string
}

func (e *kubeTaskEvent) Source() string {
	return model.KubernetesProvider
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
