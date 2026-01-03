package events

import (
	"reflect"
	"time"
)

type DomainEvent interface {
	EventType() string
	EventTime() time.Time
	EventPayload() any
}

type DomainEventEnvelope[T any] struct {
	eventType string
	occurred  time.Time
	payload   T
}

func NewDomainEventEnvelope[T any](payload T) DomainEventEnvelope[T] {
	return DomainEventEnvelope[T]{
		eventType: EventTypeFromPayload[T](),
		occurred:  time.Now(),
		payload:   payload,
	}
}

func NewDomainEventEnvelopeHydrate[T any](occured time.Time, payload T) DomainEventEnvelope[T] {
	return DomainEventEnvelope[T]{
		eventType: EventTypeFromPayload[T](),
		occurred:  occured,
		payload:   payload,
	}
}

func (e DomainEventEnvelope[T]) EventType() string    { return e.eventType }
func (e DomainEventEnvelope[T]) EventTime() time.Time { return e.occurred }
func (e DomainEventEnvelope[T]) EventPayload() any    { return e.payload }

func EventTypeFromPayload[T any]() string {
	var payload T
	t := reflect.TypeOf(payload)

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	typeName := t.Name()
	packageName := t.PkgPath()

	return packageName + "." + typeName
}
