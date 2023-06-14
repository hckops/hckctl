package event

import (
	"fmt"
	"sync"
)

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
	Kind() EventKind
	Source() string
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
