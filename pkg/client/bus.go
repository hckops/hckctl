package client

import (
	"fmt"
	"sync"
)

// TODO review: source vs provider + add visibility/type: public/private/log + debug/info/error
// TODO move "box" in client pkg

type EventSource uint8

const (
	DockerSource EventSource = iota
	KubeSource
	ArgoSource
	CloudSource
	BoxSource
)

func (e EventSource) String() string {
	return []string{"docker", "kube", "argo", "cloud", "box"}[e]
}

type Event interface {
	Source() EventSource
	fmt.Stringer
}

type EventBus struct {
	eventChan chan Event
	wg        sync.WaitGroup
}

func NewEventBus() *EventBus {
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
