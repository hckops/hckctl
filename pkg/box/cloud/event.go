package cloud

import (
	"fmt"

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

func newApiCreateCloudLoaderEvent(address string, templateName string) *cloudEvent {
	return &cloudEvent{kind: event.LoaderUpdate, value: fmt.Sprintf("loading %s/%s", address, templateName)}
}

func newApiCreateCloudEvent(templateName string, boxName string) *cloudEvent {
	return &cloudEvent{kind: event.LogDebug, value: fmt.Sprintf("api create: templateName=%s boxName=%s", templateName, boxName)}
}
