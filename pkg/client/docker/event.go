package docker

import (
	"fmt"

	"github.com/hckops/hckctl/pkg/client"
)

type dockerEventKind uint8

const (
	initClient dockerEventKind = iota
	closeClient
	setupImage
	pullImage
	createContainer
	bindPort
	execContainer
	execContainerWaiting
	execContainerError
	removeContainer
)

type DockerEvent struct {
	Kind  dockerEventKind
	Value string
}

func (e *DockerEvent) Source() client.EventSource {
	return client.DockerSource
}

func (e *DockerEvent) String() string {
	return e.Value
}

func IsDockerEvent(event client.Event) (*DockerEvent, bool) {
	if event.Source() == client.DockerSource {
		casted, ok := event.(*DockerEvent)
		return casted, ok
	}
	return nil, false
}

func newInitClientDockerEvent() *DockerEvent {
	return &DockerEvent{Kind: initClient, Value: "init docker client"}
}

func newCloseClientDockerEvent() *DockerEvent {
	return &DockerEvent{Kind: closeClient, Value: "close docker client"}
}

func newSetupImageDockerEvent(imageName string) *DockerEvent {
	return &DockerEvent{Kind: setupImage, Value: fmt.Sprintf("setup image: imageName=%s", imageName)}
}

func newPullImageDockerEvent(imageName string) *DockerEvent {
	return &DockerEvent{Kind: pullImage, Value: fmt.Sprintf("pulling %s", imageName)}
}

func newCreateContainerDockerEvent(containerName string) *DockerEvent {
	return &DockerEvent{Kind: createContainer, Value: fmt.Sprintf("create container: containerName=%s", containerName)}
}

func NewBindPortDockerEvent(message string) *DockerEvent {
	return &DockerEvent{Kind: bindPort, Value: "TODO"}
}

func newExecContainerDockerEvent(containerId string) *DockerEvent {
	return &DockerEvent{Kind: execContainer, Value: fmt.Sprintf("exec container: containerId=%s", containerId)}
}

func newExecContainerWaitingDockerEvent(containerId string) *DockerEvent {
	return &DockerEvent{Kind: execContainerWaiting, Value: fmt.Sprintf("exec container waiting: containerId=%s", containerId)}
}

func newExecContainerErrorDockerEvent(containerId string, err error) *DockerEvent {
	return &DockerEvent{Kind: execContainerError, Value: fmt.Sprintf("exec container failure: containerId=%s error=%v", containerId, err)}
}

func newRemoveContainerDockerEvent(containerId string) *DockerEvent {
	return &DockerEvent{Kind: removeContainer, Value: fmt.Sprintf("remove container: containerId=%s", containerId)}
}
