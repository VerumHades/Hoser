package outbox

import (
	"common/pkg/shared"
	"context"
	"time"

	"github.com/google/uuid"
)

// OutboxCommandRepository defines write operations for OutboxEvent.
type OutboxCommandRepository interface {
	// Create persists a new outbox event within a transaction.
	Create(
		ctx context.Context,
		event *OutboxEvent,
	) (*OutboxEvent, error)

	// MarkDispatched marks an outbox event as dispatched within a transaction.
	MarkDispatched(
		ctx context.Context,
		eventID uuid.UUID,
		dispatchedAt time.Time,
	) error

	// Delete removes an outbox event by its ID within a transaction.
	Delete(
		ctx context.Context,
		eventID uuid.UUID,
	) error
}

/*
OutboxEventCursor represents a stable pagination position
for iterating over pending events.
*/
type OutboxEventCursor struct {
	LastOccurredAt time.Time
	LastID         uuid.UUID
}

// OutboxQueryRepository defines read-only operations for OutboxEvent.
type OutboxQueryRepository interface {
	// GetByID retrieves an outbox event by its ID.
	GetByID(
		ctx context.Context,
		eventID uuid.UUID,
	) (*OutboxEvent, error)

	// FetchNextPendingBatch retrieves a batch of undispatched events after a given cursor.
	//
	// The cursor should represent the last processed event (by timestamp + ID for stability).
	FetchNextPendingBatch(
		ctx context.Context,
		request shared.BatchRequest[OutboxEventCursor],
	) (events []*OutboxEvent, nextCursor OutboxEventCursor, err error)

	// Exists checks whether an outbox event with the given ID exists.
	Exists(
		ctx context.Context,
		eventID uuid.UUID,
	) (bool, error)
}
