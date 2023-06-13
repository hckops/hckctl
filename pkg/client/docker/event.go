package docker

import (
	"fmt"

	"github.com/hckops/hckctl/pkg/client"
)

// only internal events for log purposes
type dockerEventKind uint8

const (
	clientInit dockerEventKind = iota
	clientClose
	imagePull
	imageRemove
	imageRemoveError
	containerCreate
	containerAttach
	containerAttachExit
	containerAttachError
	containerRemove
	containerList
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

func newClientInitDockerEvent() *dockerEvent {
	return &dockerEvent{kind: clientInit, value: "init docker client"}
}

func newClientCloseDockerEvent() *dockerEvent {
	return &dockerEvent{kind: clientClose, value: "close docker client"}
}

func newImagePullDockerEvent(imageName string) *dockerEvent {
	return &dockerEvent{kind: imagePull, value: fmt.Sprintf("image pull: imageName=%s", imageName)}
}

func newImageRemoveDockerEvent(imageId string) *dockerEvent {
	return &dockerEvent{kind: imageRemove, value: fmt.Sprintf("image remove: imageId=%s", imageId)}
}

// TODO level WARN
func newImageRemoveErrorDockerEvent(imageId string, err error) *dockerEvent {
	return &dockerEvent{kind: imageRemoveError, value: fmt.Sprintf("image remove error: imageId=%s error=%v", imageId, err)}
}

func newContainerCreateDockerEvent(containerName string) *dockerEvent {
	return &dockerEvent{kind: containerCreate, value: fmt.Sprintf("container create: containerName=%s", containerName)}
}

func newContainerAttachDockerEvent(containerId string) *dockerEvent {
	return &dockerEvent{kind: containerAttach, value: fmt.Sprintf("container attach: containerId=%s", containerId)}
}

func newContainerAttachExitDockerEvent(containerId string) *dockerEvent {
	return &dockerEvent{kind: containerAttachExit, value: fmt.Sprintf("container attach exit: containerId=%s", containerId)}
}

func newContainerAttachErrorDockerEvent(containerId string, err error) *dockerEvent {
	return &dockerEvent{kind: containerAttachError, value: fmt.Sprintf("container attach error: containerId=%s error=%v", containerId, err)}
}

func newContainerRemoveDockerEvent(containerId string) *dockerEvent {
	return &dockerEvent{kind: containerRemove, value: fmt.Sprintf("container remove: containerId=%s", containerId)}
}

func newContainerListDockerEvent(index int, containerId string, containerName string) *dockerEvent {
	return &dockerEvent{kind: containerList, value: fmt.Sprintf("container list: (%d) containerId=%s containerName=%s", index, containerId, containerName)}
}
