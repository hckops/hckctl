package kubernetes

import (
	"github.com/hckops/hckctl/pkg/event"
)

type kubeEvent struct {
	kind  event.EventKind
	value string
}

func (e *kubeEvent) Source() string {
	return "kube"
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
