package box

import (
	"fmt"
	"github.com/hckops/hckctl/pkg/client"
	"github.com/hckops/hckctl/pkg/template/model"
)

type boxEventKind uint8

const (
	Debug boxEventKind = iota
	Console
	LoaderUpdate
	LoaderClose
)

type BoxEvent struct {
	Kind  boxEventKind
	value string
}

func (e *BoxEvent) Source() client.EventSource {
	return client.BoxSource
}

func (e *BoxEvent) String() string {
	return e.value
}

func IsBoxEvent(event client.Event) (*BoxEvent, bool) {
	if event.Source() == client.BoxSource {
		casted, ok := event.(*BoxEvent)
		return casted, ok
	}
	return nil, false
}

func newGenericBoxEvent(message string, values ...any) *BoxEvent {
	return &BoxEvent{Kind: Debug, value: fmt.Sprintf(message, values...)}
}

func newInitBoxEvent() *BoxEvent {
	return &BoxEvent{Kind: Debug, value: "init box client"}
}

func newBindPortBoxEvent(boxName string, port model.BoxPort) *BoxEvent {
	return &BoxEvent{Kind: Console, value: fmt.Sprintf(
		"[%s][%s]   \texpose (remote) %s -> (local) http://localhost:%s",
		boxName, port.Alias, port.Remote, port.Local)}
}

func newPullImageBoxEvent(imageName string) *BoxEvent {
	return &BoxEvent{Kind: LoaderUpdate, value: fmt.Sprintf("pulling image %s", imageName)}
}

func newContainerWaitingBoxEvent() *BoxEvent {
	return &BoxEvent{Kind: LoaderClose, value: "waiting"}
}
