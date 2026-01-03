package outbox

import (
	"time"

	"github.com/google/uuid"
)

/*
AggregateType represents the type of aggregate that produced an event.
*/
type AggregateType string

const (
	AggregateTypeListing AggregateType = "listing"
)

/*
EventType represents the semantic type of a domain event.
*/
type EventType string

const (
	EventTypeCreate EventType = "created"
	EventTypeUpdate EventType = "updated"
	EventTypeDelete EventType = "deleted"
)

/*
OutboxEvent represents a transactional record of a domain event
that must be delivered to external systems. It is immutable
except for dispatch metadata.
*/
type OutboxEvent struct {
	id            uuid.UUID
	aggregateType AggregateType
	aggregateID   string
	eventType     EventType
	payload       []byte
	occurredAt    time.Time
}

/*
NewOutboxEvent creates a new undispatched outbox event.
*/
func NewOutboxEvent(
	aggregateType AggregateType,
	aggregateID string,
	eventType EventType,
	payload []byte,
	occurredAt time.Time,
) OutboxEvent {
	return OutboxEvent{
		id:            uuid.New(),
		aggregateType: aggregateType,
		aggregateID:   aggregateID,
		eventType:     eventType,
		payload:       payload,
		occurredAt:    occurredAt.UTC(),
	}
}

func (event OutboxEvent) ID() uuid.UUID {
	return event.id
}

func (event OutboxEvent) AggregateType() AggregateType {
	return event.aggregateType
}

func (event OutboxEvent) AggregateID() string {
	return event.aggregateID
}

func (event OutboxEvent) EventType() EventType {
	return event.eventType
}

func (event OutboxEvent) Payload() []byte {
	return event.payload
}

func (event OutboxEvent) OccurredAt() time.Time {
	return event.occurredAt
}

func ptrTime(value time.Time) *time.Time {
	return &value
}
