package cloud

import (
	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/event"
)

type cloudEvent struct {
	kind  event.EventKind
	value string
}

func (e *cloudEvent) Source() string {
	return model.Cloud.String()
}

func (e *cloudEvent) Kind() event.EventKind {
	return e.kind
}

func (e *cloudEvent) String() string {
	return e.value
}

func newClientInitCloudEvent() *cloudEvent {
	return &cloudEvent{kind: event.LogDebug, value: "init cloud client"}
}

func newClientCloseCloudEvent() *cloudEvent {
	return &cloudEvent{kind: event.LogDebug, value: "close cloud client"}
}
