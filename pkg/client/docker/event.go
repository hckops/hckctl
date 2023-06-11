package docker

import (
	"fmt"

	"github.com/hckops/hckctl/pkg/client"
)

// only internal events for log purposes
type dockerEventKind uint8

const (
	initClient dockerEventKind = iota
	closeClient
	setupImage
	pullImage
	createContainer
	execContainer
	execContainerWaiting
	execContainerError
	execContainerExit
	removeContainer
	listContainers
)

type dockerEvent struct {
	kind  dockerEventKind
	value string
}

func (e *dockerEvent) Source() client.EventSource {
	return client.DockerSource
}

func (e *dockerEvent) String() string {
	return e.value
}

func newInitClientDockerEvent() *dockerEvent {
	return &dockerEvent{kind: initClient, value: "init docker client"}
}

func newCloseClientDockerEvent() *dockerEvent {
	return &dockerEvent{kind: closeClient, value: "close docker client"}
}

func newSetupImageDockerEvent(imageName string) *dockerEvent {
	return &dockerEvent{kind: setupImage, value: fmt.Sprintf("setup image: imageName=%s", imageName)}
}

func newPullImageDockerEvent(imageName string) *dockerEvent {
	return &dockerEvent{kind: pullImage, value: fmt.Sprintf("pull image %s", imageName)}
}

func newCreateContainerDockerEvent(containerName string) *dockerEvent {
	return &dockerEvent{kind: createContainer, value: fmt.Sprintf("create container: containerName=%s", containerName)}
}

func newExecContainerDockerEvent(containerId string) *dockerEvent {
	return &dockerEvent{kind: execContainer, value: fmt.Sprintf("exec container: containerId=%s", containerId)}
}

func newExecContainerWaitingDockerEvent(containerId string) *dockerEvent {
	return &dockerEvent{kind: execContainerWaiting, value: fmt.Sprintf("exec container waiting: containerId=%s", containerId)}
}

func newExecContainerErrorDockerEvent(containerId string, err error) *dockerEvent {
	return &dockerEvent{kind: execContainerError, value: fmt.Sprintf("exec container failure: containerId=%s error=%v", containerId, err)}
}

func newExecContainerExitDockerEvent(containerId string) *dockerEvent {
	return &dockerEvent{kind: execContainerExit, value: fmt.Sprintf("exec container exit: containerId=%s", containerId)}
}

func newRemoveContainerDockerEvent(containerId string) *dockerEvent {
	return &dockerEvent{kind: removeContainer, value: fmt.Sprintf("remove container: containerId=%s", containerId)}
}

func newListContainersDockerEvent(index int, containerId string, containerName string) *dockerEvent {
	return &dockerEvent{kind: listContainers, value: fmt.Sprintf("[%d] list containers: containerId=%s containerName=%s", index, containerId, containerName)}
}
