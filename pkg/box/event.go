package box

import (
	"fmt"
	"sync"
)

// TODO move in a different pkg
// https://go.dev/blog/pipelines

type EventKind uint8

const (
	DebugEvent EventKind = iota
	InfoEvent
	ErrorEvent
	PriorityEvent
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

func (bus *EventBus) PublishDebugEvent(source, message string, values ...any) {
	bus.publishEvent(DebugEvent, source, fmt.Sprintf(message, values...))
}

func (bus *EventBus) PublishInfoEvent(source, message string, values ...any) {
	bus.publishEvent(InfoEvent, source, fmt.Sprintf(message, values...))
}

func (bus *EventBus) PublishErrorEvent(source, message string, values ...any) {
	bus.publishEvent(ErrorEvent, source, fmt.Sprintf(message, values...))
}

func (bus *EventBus) PublishPriorityEvent(source string, message string, values ...any) {
	bus.publishEvent(PriorityEvent, source, fmt.Sprintf(message, values...))
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

func (bus *EventBus) Drain() {
	bus.SubscribeEvents(func(event Event) {})
}

func (bus *EventBus) Close() {
	bus.wg.Wait()
}
