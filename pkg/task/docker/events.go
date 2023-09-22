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

func newInitDockerClientEvent() *dockerTaskEvent {
	return &dockerTaskEvent{kind: event.LogDebug, value: "init docker client"}
}

func newCloseDockerClientEvent() *dockerTaskEvent {
	return &dockerTaskEvent{kind: event.LogDebug, value: "close docker client"}
}

func newImagePullDockerEvent(imageName string) *dockerTaskEvent {
	return &dockerTaskEvent{kind: event.LogInfo, value: fmt.Sprintf("image pull: imageName=%s", imageName)}
}

func newImagePullDockerLoaderEvent(imageName string) *dockerTaskEvent {
	return &dockerTaskEvent{kind: event.LoaderUpdate, value: fmt.Sprintf("pulling image %s", imageName)}
}

func newImagePullIgnoreDockerEvent(imageName string) *dockerTaskEvent {
	return &dockerTaskEvent{kind: event.LogWarning, value: fmt.Sprintf("image pull ignored: imageName=%s", imageName)}
}

func newImageRemoveDockerEvent(imageId string) *dockerTaskEvent {
	return &dockerTaskEvent{kind: event.LogInfo, value: fmt.Sprintf("image remove: imageId=%s", imageId)}
}

func newImageRemoveIgnoreDockerEvent(imageId string, err error) *dockerTaskEvent {
	return &dockerTaskEvent{kind: event.LogWarning, value: fmt.Sprintf("image remove ignored: imageId=%s error=%v", imageId, err)}
}

func newNetworkUpsertDockerEvent(networkName string, networkId string) *dockerTaskEvent {
	return &dockerTaskEvent{kind: event.LogInfo, value: fmt.Sprintf("network upsert: networkName=%s networkId=%s", networkName, networkId)}
}

func newContainerCreateStatusDockerEvent(status string) *dockerTaskEvent {
	return &dockerTaskEvent{kind: event.LogDebug, value: status}
}

func newContainerCreateDockerEvent(templateName string, containerName string, containerId string) *dockerTaskEvent {
	return &dockerTaskEvent{kind: event.LogInfo, value: fmt.Sprintf("container create: templateName=%s containerName=%s containerId=%s", templateName, containerName, containerId)}
}

func newContainerCreateDockerLoaderEvent() *dockerTaskEvent {
	return &dockerTaskEvent{kind: event.LoaderUpdate, value: "running"}
}

func newContainerStartDockerLoaderEvent() *dockerTaskEvent {
	return &dockerTaskEvent{kind: event.LoaderStop, value: "waiting"}
}

func newVpnConnectDockerLoaderEvent(vpnName string) *dockerTaskEvent {
	return &dockerTaskEvent{kind: event.LoaderUpdate, value: fmt.Sprintf("connecting to %s", vpnName)}
}
