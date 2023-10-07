package cloud

import (
	"fmt"

	"github.com/hckops/hckctl/pkg/event"
	"github.com/hckops/hckctl/pkg/lab/model"
)

type cloudLabEvent struct {
	kind  event.EventKind
	value string
}

func (e *cloudLabEvent) Source() string {
	return model.Cloud.String()
}

func (e *cloudLabEvent) Kind() event.EventKind {
	return e.kind
}

func (e *cloudLabEvent) String() string {
	return e.value
}

func newInitCloudClientEvent() *cloudLabEvent {
	return &cloudLabEvent{kind: event.LogDebug, value: "init cloud client"}
}

func newApiCreateCloudLoaderEvent(address string, templateName string) *cloudLabEvent {
	return &cloudLabEvent{kind: event.LoaderUpdate, value: fmt.Sprintf("loading %s/%s", address, templateName)}
}

func newApiCreateCloudEvent(templateName string, labName string) *cloudLabEvent {
	return &cloudLabEvent{kind: event.LogInfo, value: fmt.Sprintf("api create: templateName=%s labName=%s", templateName, labName)}
}
