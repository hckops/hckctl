package box

import (
	"sync"
)

// TODO move in a different pkg
// https://go.dev/blog/pipelines

type EventKind uint8

const (
	DebugEvent EventKind = iota
	InfoEvent
	SuccessEvent
	ErrorEvent
)

type Event struct {
	Kind    EventKind
	Source  string
	Message string
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

func (bus *EventBus) publishEvent(kind EventKind, source, message string) {
	bus.wg.Add(1)
	go func() {
		bus.eventChan <- Event{
			Kind:    kind,
			Source:  source,
			Message: message,
		}
		bus.wg.Done()
	}()
}

func (bus *EventBus) Close() {
	bus.wg.Wait()
}

func (bus *EventBus) PublishDebugEvent(source, message string) {
	bus.publishEvent(DebugEvent, source, message)
}

func (bus *EventBus) PublishEmptySuccessEvent(source string) {
	bus.publishEvent(SuccessEvent, source, "")
}

func (bus *EventBus) SubscribeEvents(callback func(event Event)) {
	go func() {
		for {
			select {
			case event := <-bus.eventChan:
				callback(event)
			}
		}
	}()
}
