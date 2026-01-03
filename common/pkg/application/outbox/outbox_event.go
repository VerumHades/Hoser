package outbox

import (
	"common/pkg/shared"
	"time"
)

/*
OutboxEvent represents a transactional record of a domain event
that must be delivered to external systems. It is immutable
except for dispatch metadata.
*/
type OutboxEvent struct {
	id         shared.OutboxEventID
	eventType  string
	payload    []byte
	occurredAt time.Time
}

/*
NewOutboxEvent creates a new undispatched outbox event.
*/
func NewOutboxEvent(
	eventType string,
	payload []byte,
	occurredAt time.Time,
) *OutboxEvent {
	return NewOutboxEventWithID(
		shared.OutboxEventID(shared.GenerateUUID()),
		eventType,
		payload,
		occurredAt.UTC(),
	)
}
func NewOutboxEventWithID(
	id shared.OutboxEventID,
	eventType string,
	payload []byte,
	occurredAt time.Time,
) *OutboxEvent {
	return &OutboxEvent{
		id:         id,
		eventType:  eventType,
		payload:    payload,
		occurredAt: occurredAt.UTC(),
	}
}

func (event OutboxEvent) ID() shared.OutboxEventID {
	return event.id
}

func (event OutboxEvent) EventType() string {
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
