package docker

import (
	"fmt"

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

func newImagePullDockerLoaderEvent(imageName string) *dockerTaskEvent {
	return &dockerTaskEvent{kind: event.LoaderUpdate, value: fmt.Sprintf("pulling image %s", imageName)}
}

func newNetworkUpsertDockerEvent(networkName string, networkId string) *dockerTaskEvent {
	return &dockerTaskEvent{kind: event.LogInfo, value: fmt.Sprintf("network upsert: networkName=%s networkId=%s", networkName, networkId)}
}

func newVolumeMountDockerEvent(containerId string, hostDir string, containerDir string) *dockerTaskEvent {
	return &dockerTaskEvent{kind: event.LogInfo, value: fmt.Sprintf("volume mount: containerId=%s hostDir=%s containerDir=%s", containerId, hostDir, containerDir)}
}

func newContainerCreateStatusDockerEvent(status string) *dockerTaskEvent {
	return &dockerTaskEvent{kind: event.LogDebug, value: status}
}

func newContainerCreateDockerEvent(templateName string, containerName string, containerId string) *dockerTaskEvent {
	return &dockerTaskEvent{kind: event.LogInfo, value: fmt.Sprintf("container create: templateName=%s containerName=%s containerId=%s", templateName, containerName, containerId)}
}

func newContainerLogDockerEvent(logFileName string) *dockerTaskEvent {
	return &dockerTaskEvent{kind: event.LogInfo, value: fmt.Sprintf("container log: logFileName=%s", logFileName)}
}

func newContainerLogDockerConsoleEvent(logFileName string) *dockerTaskEvent {
	return &dockerTaskEvent{kind: event.PrintConsole, value: fmt.Sprintf("output file: %s", logFileName)}
}

func newContainerCreateDockerLoaderEvent() *dockerTaskEvent {
	return &dockerTaskEvent{kind: event.LoaderUpdate, value: "running"}
}

func newContainerWaitDockerLoaderEvent() *dockerTaskEvent {
	return &dockerTaskEvent{kind: event.LoaderStop, value: "waiting"}
}
