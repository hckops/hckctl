package docker

import (
	"fmt"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/event"
)

type dockerBoxEvent struct {
	kind  event.EventKind
	value string
}

func (e *dockerBoxEvent) Source() string {
	return model.Docker.String()
}

func (e *dockerBoxEvent) Kind() event.EventKind {
	return e.kind
}

func (e *dockerBoxEvent) String() string {
	return e.value
}

func newImagePullDockerLoaderEvent(imageName string) *dockerBoxEvent {
	return &dockerBoxEvent{kind: event.LoaderUpdate, value: fmt.Sprintf("pulling image %s", imageName)}
}

func newNetworkUpsertDockerEvent(networkName string, networkId string) *dockerBoxEvent {
	return &dockerBoxEvent{kind: event.LogInfo, value: fmt.Sprintf("network upsert: networkName=%s networkId=%s", networkName, networkId)}
}

func newContainerCreatePortBindDockerEvent(containerName string, port model.BoxPort) *dockerBoxEvent {
	return &dockerBoxEvent{kind: event.LogInfo, value: fmt.Sprintf(
		"container create port bind: containerName=%s portAlias=%s portRemote=%s portLocal=%s",
		containerName, port.Alias, port.Remote, port.Local)}
}

func newContainerCreatePortBindDockerConsoleEvent(containerName string, port model.BoxPort, padding int) *dockerBoxEvent {
	return &dockerBoxEvent{kind: event.PrintConsole, value: fmt.Sprintf(
		"[%s][%-*s] tunnel (remote) %s -> (local) %s",
		containerName, padding, port.Alias, port.Remote, port.Local)}
}

func newContainerCreateEnvDockerEvent(containerName string, env model.BoxEnv) *dockerBoxEvent {
	return &dockerBoxEvent{kind: event.LogInfo, value: fmt.Sprintf("container create env: containerName=%s key=%s value=%s", containerName, env.Key, env.Value)}
}

func newContainerCreateEnvDockerConsoleEvent(containerName string, env model.BoxEnv) *dockerBoxEvent {
	return &dockerBoxEvent{kind: event.PrintConsole, value: fmt.Sprintf("[%s] %s=%s", containerName, env.Key, env.Value)}
}

func newContainerCreateStatusDockerEvent(status string) *dockerBoxEvent {
	return &dockerBoxEvent{kind: event.LogDebug, value: status}
}

func newContainerCreateDockerEvent(templateName string, containerName string, containerId string) *dockerBoxEvent {
	return &dockerBoxEvent{kind: event.LogInfo, value: fmt.Sprintf("container create: templateName=%s containerName=%s containerId=%s", templateName, containerName, containerId)}
}

func newContainerRestartDockerEvent(containerId string, status string) *dockerBoxEvent {
	return &dockerBoxEvent{kind: event.LogInfo, value: fmt.Sprintf("container restart: containerId=%s status=%s", containerId, status)}
}

func newContainerExecDockerEvent(containerName string, containerId string, command string) *dockerBoxEvent {
	return &dockerBoxEvent{kind: event.LogInfo, value: fmt.Sprintf("container exec: containerName=%s containerId=%s command=%s", containerName, containerId, command)}
}

func newContainerExecIgnoreDockerEvent(containerId string) *dockerBoxEvent {
	return &dockerBoxEvent{kind: event.LogWarning, value: fmt.Sprintf("container exec connection ignored: containerId=%s", containerId)}
}

func newContainerExecDockerLoaderEvent() *dockerBoxEvent {
	return &dockerBoxEvent{kind: event.LoaderStop, value: "waiting"}
}

func newContainerExecExitDockerEvent(containerId string) *dockerBoxEvent {
	return &dockerBoxEvent{kind: event.LogDebug, value: fmt.Sprintf("container exec exit: containerId=%s", containerId)}
}

func newContainerExecErrorDockerEvent(containerId string, err error) *dockerBoxEvent {
	return &dockerBoxEvent{kind: event.LogError, value: fmt.Sprintf("container exec error: containerId=%s error=%v", containerId, err)}
}

func newContainerLogsDockerEvent(containerId string) *dockerBoxEvent {
	return &dockerBoxEvent{kind: event.LogDebug, value: fmt.Sprintf("container logs: containerId=%s", containerId)}
}

func newContainerLogsExitDockerEvent(containerId string) *dockerBoxEvent {
	return &dockerBoxEvent{kind: event.LogDebug, value: fmt.Sprintf("container logs exit: containerId=%s", containerId)}
}

func newContainerLogsErrorDockerEvent(containerId string, err error) *dockerBoxEvent {
	return &dockerBoxEvent{kind: event.LogError, value: fmt.Sprintf("container logs error: containerId=%s error=%v", containerId, err)}
}

func newContainerListDockerEvent(index int, containerName string, containerId string, healthy bool) *dockerBoxEvent {
	return &dockerBoxEvent{kind: event.LogDebug, value: fmt.Sprintf("container list: (%d) containerName=%s containerId=%s healthy=%v", index, containerName, containerId, healthy)}
}

func newContainerRemoveDockerEvent(containerName string, containerId string) *dockerBoxEvent {
	return &dockerBoxEvent{kind: event.LogInfo, value: fmt.Sprintf("container remove: containerName=%s containerId=%s", containerName, containerId)}
}

func newContainerRemoveIgnoreDockerEvent(containerName string, containerId string, err error) *dockerBoxEvent {
	return &dockerBoxEvent{kind: event.LogWarning, value: fmt.Sprintf("container remove ignored: containerName=%s containerId=%s error=%v", containerName, containerId, err)}
}

func newContainerInspectDockerEvent(containerId string) *dockerBoxEvent {
	return &dockerBoxEvent{kind: event.LogInfo, value: fmt.Sprintf("container inspect: containerId=%s", containerId)}
}
