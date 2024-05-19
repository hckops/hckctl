package docker

import (
	"fmt"

	"github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/event"
)

type dockerCommonEvent struct {
	kind  event.EventKind
	value string
}

func (e *dockerCommonEvent) Source() string {
	return model.DockerProvider
}

func (e *dockerCommonEvent) Kind() event.EventKind {
	return e.kind
}

func (e *dockerCommonEvent) String() string {
	return e.value
}

func newInitDockerClientEvent() *dockerCommonEvent {
	return &dockerCommonEvent{kind: event.LogDebug, value: "init docker client"}
}

func newCloseDockerClientEvent() *dockerCommonEvent {
	return &dockerCommonEvent{kind: event.LogDebug, value: "close docker client"}
}

func newImagePullDockerEvent(imageName string, platform string) *dockerCommonEvent {
	return &dockerCommonEvent{kind: event.LogInfo, value: fmt.Sprintf("image pull: imageName=%s platform=%s", imageName, platform)}
}

func newImagePullIgnoreDockerEvent(imageName string) *dockerCommonEvent {
	return &dockerCommonEvent{kind: event.LogWarning, value: fmt.Sprintf("image pull ignored: imageName=%s", imageName)}
}

func newImageRemoveDockerEvent(imageId string) *dockerCommonEvent {
	return &dockerCommonEvent{kind: event.LogInfo, value: fmt.Sprintf("image remove: imageId=%s", imageId)}
}

func newImageRemoveIgnoreDockerEvent(imageId string, err error) *dockerCommonEvent {
	return &dockerCommonEvent{kind: event.LogWarning, value: fmt.Sprintf("image remove ignored: imageId=%s error=%v", imageId, err)}
}

func newSidecarVpnCreateDockerEvent(containerName string, containerId string) *dockerCommonEvent {
	return &dockerCommonEvent{kind: event.LogInfo, value: fmt.Sprintf("sidecar-vpn create: containerName=%s containerId=%s", containerName, containerId)}
}

func newSidecarVpnCreateStatusDockerEvent(status string) *dockerCommonEvent {
	return &dockerCommonEvent{kind: event.LogDebug, value: status}
}

func newSidecarVpnConnectDockerEvent(vpnName string) *dockerCommonEvent {
	return &dockerCommonEvent{kind: event.LogInfo, value: fmt.Sprintf("sidecar-vpn connect: vpnName=%s", vpnName)}
}

func newSidecarVpnConnectDockerLoaderEvent(vpnName string) *dockerCommonEvent {
	return &dockerCommonEvent{kind: event.LoaderUpdate, value: fmt.Sprintf("connecting to %s", vpnName)}
}
