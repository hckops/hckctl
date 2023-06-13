package box

import (
	"fmt"
	"sync"
)

type EventSource uint8

const (
	DockerSource EventSource = iota
	KubeSource
	ArgoSource
	CloudSource
)

func (e EventSource) String() string {
	return []string{"docker", "kube", "argo", "cloud", "box"}[e]
}

type EventKind uint8

const (
	LogDebug EventKind = iota
	LogInfo
	LogWarning
	LogError
	PrintConsole
	LoaderUpdate
	LoaderStop
)

func (e EventKind) String() string {
	return []string{"debug", "info", "warning", "console", "update", "stop"}[e]
}

type Event interface {
	Source() EventSource
	Kind() EventKind
	fmt.Stringer
}

type EventBus struct {
	eventChan chan Event
	wg        sync.WaitGroup
}

func newEventBus() *EventBus {
	return &EventBus{
		eventChan: make(chan Event),
	}
}

func (bus *EventBus) Publish(event Event) {
	bus.wg.Add(1)
	go func() {
		bus.eventChan <- event
		bus.wg.Done()
	}()
}

func (bus *EventBus) Subscribe(callback func(event Event)) {
	go func() {
		for {
			select {
			case event := <-bus.eventChan:
				callback(event)
			}
		}
	}()
}

func (bus *EventBus) Drain() {
	bus.Subscribe(func(event Event) {})
}

func (bus *EventBus) Close() {
	bus.wg.Wait()
}
