package docker

import (
	"github.com/hckops/hckctl/pkg/client"
)

type DockerEventKind uint8

const (
	Create DockerEventKind = iota
	Open
	Remove
	Exec
)

type DockerEvent struct {
	Kind  DockerEventKind
	Value string
}

func (e *DockerEvent) Source() client.EventSource {
	return client.DockerSource
}

func (e *DockerEvent) String() string {
	return "todo-string"
}

func IsDockerEvent(event client.Event) (*DockerEvent, bool) {
	if event.Source() == client.DockerSource {
		casted, ok := event.(*DockerEvent)
		return casted, ok
	}
	return nil, false
}

func NewDockerCreateEvent(value string) *DockerEvent {
	return &DockerEvent{Kind: Create, Value: value}
}
