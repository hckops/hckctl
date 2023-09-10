package docker

import (
	"github.com/hckops/hckctl/pkg/event"
	"github.com/hckops/hckctl/pkg/task/model"
)

type dockerTaskEvent struct {
	kind  event.EventKind
	value string
}

func (e *dockerTaskEvent) Source() string {
	return model.Docker.String()
}

func (e *dockerTaskEvent) Kind() event.EventKind {
	return e.kind
}

func (e *dockerTaskEvent) String() string {
	return e.value
}

func newInitDockerClientEvent() *dockerTaskEvent {
	return &dockerTaskEvent{kind: event.LogDebug, value: "init docker client"}
}

func newCloseDockerClientEvent() *dockerTaskEvent {
	return &dockerTaskEvent{kind: event.LogDebug, value: "close docker client"}
}
