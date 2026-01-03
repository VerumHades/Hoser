package outbox

import (
	"common/pkg/domain/entities/events"
	"context"
	"encoding/json"
)

// EventPublisher defines the interface for publishing domain event envelopes
// by writing them into the
type EventPublisher interface {
	PublishEnvelope(ctx context.Context, envelope events.DomainEvent) error
	PublishEnvelopes(ctx context.Context, envelopes ...events.DomainEvent) error
}

// OutboxEventPublisher implements EventPublisher on top of OutboxCommandRepository.
type OutboxEventPublisher struct {
	repo OutboxCommandRepository
}

// NewOutboxEventPublisher constructs a new OutboxEventPublisher.
func NewOutboxEventPublisher(repo OutboxCommandRepository) *OutboxEventPublisher {
	return &OutboxEventPublisher{repo: repo}
}

func (p *OutboxEventPublisher) PublishEnvelope(ctx context.Context, envelope events.DomainEvent) error {
	payload, err := json.Marshal(envelope.EventPayload())
	if err != nil {
		return err
	}
	event := NewOutboxEvent(
		envelope.EventType(),
		payload,
		envelope.EventTime(),
	)

	return p.repo.Create(ctx, event)
}
func (p *OutboxEventPublisher) PublishEnvelopes(ctx context.Context, envelopes ...events.DomainEvent) error {
	events := make([]*OutboxEvent, len(envelopes))
	for i, envelope := range envelopes {
		payload, err := json.Marshal(envelope.EventPayload())
		if err != nil {
			return err
		}
		events[i] = NewOutboxEvent(
			envelope.EventType(),
			payload,
			envelope.EventTime(),
		)
	}
	return p.repo.BulkCreate(ctx, events)
}
