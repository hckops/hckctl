package cloud

import (
	"fmt"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/event"
)

type cloudEvent struct {
	kind  event.EventKind
	value string
}

func (e *cloudEvent) Source() string {
	return model.Cloud.String()
}

func (e *cloudEvent) Kind() event.EventKind {
	return e.kind
}

func (e *cloudEvent) String() string {
	return e.value
}

func newInitCloudClientEvent() *cloudEvent {
	return &cloudEvent{kind: event.LogDebug, value: "init cloud client"}
}

func newCloseCloudClientEvent() *cloudEvent {
	return &cloudEvent{kind: event.LogDebug, value: "close cloud client"}
}

func newApiRawCloudEvent(value string) *cloudEvent {
	return &cloudEvent{kind: event.LogDebug, value: value}
}

func newApiCreateCloudLoaderEvent(address string, templateName string) *cloudEvent {
	return &cloudEvent{kind: event.LoaderUpdate, value: fmt.Sprintf("loading %s/%s", address, templateName)}
}

func newApiCreateCloudEvent(templateName string, boxName string, size string) *cloudEvent {
	return &cloudEvent{kind: event.LogInfo, value: fmt.Sprintf("api create: templateName=%s boxName=%s size=%s", templateName, boxName, size)}
}

func newApiExecCloudEvent(boxName string) *cloudEvent {
	return &cloudEvent{kind: event.LogInfo, value: fmt.Sprintf("api exec: boxName=%s", boxName)}
}

func newApiExecErrorCloudEvent(boxName string, err error) *cloudEvent {
	return &cloudEvent{kind: event.LogError, value: fmt.Sprintf("api exec error: boxName=%s error=%v", boxName, err)}
}

func newApiStopCloudLoaderEvent() *cloudEvent {
	return &cloudEvent{kind: event.LoaderStop, value: "waiting"}
}

func newApiEnvCloudEvent(boxName string, env model.BoxEnv) *cloudEvent {
	return &cloudEvent{kind: event.LogInfo, value: fmt.Sprintf("api env: boxName=%s key=%s value=%s", boxName, env.Key, env.Value)}
}

func newApiEnvCloudConsoleEvent(boxName string, env model.BoxEnv) *cloudEvent {
	return &cloudEvent{kind: event.PrintConsole, value: fmt.Sprintf("[%s] %s=%s", boxName, env.Key, env.Value)}
}

func newApiTunnelIgnoreCloudEvent(boxName string) *cloudEvent {
	return &cloudEvent{kind: event.LogWarning, value: fmt.Sprintf("api tunnel ignored: boxName=%s", boxName)}
}

func newApiTunnelBindingCloudEvent(boxName string, port model.BoxPort) *cloudEvent {
	return &cloudEvent{kind: event.LogInfo, value: fmt.Sprintf(
		"api tunnel: boxName=%s portAlias=%s portRemote=%s portLocal=%s",
		boxName, port.Alias, port.Remote, port.Local)}
}

func newApiTunnelBindingCloudConsoleEvent(boxName string, port model.BoxPort, padding int) *cloudEvent {
	return &cloudEvent{kind: event.PrintConsole, value: fmt.Sprintf(
		"[%s][%-*s] tunnel (remote) %s -> (local) %s",
		boxName, padding, port.Alias, port.Remote, port.Local)}
}

func newApiTunnelStartCloudEvent(boxName string, port model.BoxPort, connection string) *cloudEvent {
	return &cloudEvent{kind: event.LogDebug, value: fmt.Sprintf("api tunnel start: boxName=%s portRemote=%s portLocal=%s connection=%s", boxName, port.Remote, port.Local, connection)}
}

func newApiTunnelStopCloudEvent(boxName string, port model.BoxPort, connection string) *cloudEvent {
	return &cloudEvent{kind: event.LogDebug, value: fmt.Sprintf("api tunnel stop: boxName=%s portRemote=%s portLocal=%s connection=%s", boxName, port.Remote, port.Local, connection)}
}

func newApiTunnelErrorCloudEvent(boxName string, err error) *cloudEvent {
	return &cloudEvent{kind: event.LogError, value: fmt.Sprintf("api tunnel error: boxName=%s error=%v", boxName, err)}
}

func newApiTunnelListenCloudLoaderEvent() *cloudEvent {
	return &cloudEvent{kind: event.LoaderUpdate, value: "listening"}
}

func newApiDescribeCloudEvent(boxName string) *cloudEvent {
	return &cloudEvent{kind: event.LogInfo, value: fmt.Sprintf("api describe: boxName=%s", boxName)}
}

func newApiListCloudEvent(index int, boxName string) *cloudEvent {
	return &cloudEvent{kind: event.LogInfo, value: fmt.Sprintf("api list: (%d) boxName=%s", index, boxName)}
}

func newApiDeleteCloudEvent(index int, boxName string) *cloudEvent {
	return &cloudEvent{kind: event.LogInfo, value: fmt.Sprintf("api delete: (%d) boxName=%s", index, boxName)}
}
