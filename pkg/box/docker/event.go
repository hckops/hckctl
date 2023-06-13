package docker

import (
	"fmt"

	"github.com/hckops/hckctl/pkg/box"
	"github.com/hckops/hckctl/pkg/template/model"
)

type dockerEvent struct {
	kind  box.EventKind
	value string
}

func (e *dockerEvent) Source() box.EventSource {
	return box.DockerSource
}

func (e *dockerEvent) Kind() box.EventKind {
	return e.kind
}

func (e *dockerEvent) String() string {
	return e.value
}

func newClientInitDockerEvent() *dockerEvent {
	return &dockerEvent{kind: box.LogDebug, value: "init docker client"}
}

func newClientCloseDockerEvent() *dockerEvent {
	return &dockerEvent{kind: box.LogDebug, value: "close docker client"}
}

func newImagePullDockerEvent(imageName string) *dockerEvent {
	return &dockerEvent{kind: box.LogDebug, value: fmt.Sprintf("image pull: imageName=%s", imageName)}
}

func newImagePullDockerLoaderEvent(imageName string) *dockerEvent {
	return &dockerEvent{kind: box.LoaderUpdate, value: fmt.Sprintf("pulling image %s", imageName)}
}

func newImageRemoveDockerEvent(imageId string) *dockerEvent {
	return &dockerEvent{kind: box.LogDebug, value: fmt.Sprintf("image remove: imageId=%s", imageId)}
}

func newImageRemoveErrorDockerEvent(imageId string, err error) *dockerEvent {
	return &dockerEvent{kind: box.LogWarning, value: fmt.Sprintf("image remove error: imageId=%s error=%v", imageId, err)}
}

func newContainerCreateSkipVirtualPortDockerEvent(containerName string, port model.BoxPort) *dockerEvent {
	return &dockerEvent{kind: box.LogDebug, value: fmt.Sprintf("container create skipping virtual port: containerName=%s portAlias=%s", containerName, port.Alias)}
}

func newContainerCreatePortBindDockerConsoleEvent(containerName string, port model.BoxPort) *dockerEvent {
	return &dockerEvent{kind: box.PrintConsole, value: fmt.Sprintf(
		"[%s][%s]   \texpose (remote) %s -> (local) %s",
		containerName, port.Alias, port.Remote, port.Local)}
}

func newContainerCreatePortBindDockerEvent(containerName string, port model.BoxPort) *dockerEvent {
	return &dockerEvent{kind: box.LogDebug, value: fmt.Sprintf(
		"container create port bind: containerName=%s portAlias=%s portRemote=%s portLocal=%s",
		containerName, port.Alias, port.Remote, port.Local)}
}

func newContainerCreateDockerEvent(templateName string, containerName string, containerId string) *dockerEvent {
	return &dockerEvent{kind: box.LogDebug, value: fmt.Sprintf("container create: templateName-%s containerName=%s containerId=%s", templateName, containerName, containerId)}
}

func newContainerAttachDockerEvent(containerName string, containerId string, command string) *dockerEvent {
	return &dockerEvent{kind: box.LogDebug, value: fmt.Sprintf("container attach: containerName=%s containerId=%s command=%s", containerName, containerId, command)}
}

func newContainerAttachDockerLoaderEvent() *dockerEvent {
	return &dockerEvent{kind: box.LoaderStop, value: "waiting"}
}

func newContainerAttachExitDockerEvent(containerId string) *dockerEvent {
	return &dockerEvent{kind: box.LogDebug, value: fmt.Sprintf("container attach exit: containerId=%s", containerId)}
}

func newContainerAttachErrorDockerEvent(containerId string, err error) *dockerEvent {
	return &dockerEvent{kind: box.LogDebug, value: fmt.Sprintf("container attach error: containerId=%s error=%v", containerId, err)}
}

func newContainerListDockerEvent(index int, containerName string, containerId string) *dockerEvent {
	return &dockerEvent{kind: box.LogDebug, value: fmt.Sprintf("container list: (%d) containerName=%s containerId=%s", index, containerName, containerId)}
}

func newContainerRemoveDockerEvent(containerId string) *dockerEvent {
	return &dockerEvent{kind: box.LogDebug, value: fmt.Sprintf("container remove: containerId=%s", containerId)}
}
