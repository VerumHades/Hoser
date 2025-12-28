package instance

import (
	"common/internal/shared"
	"context"
	"time"
)

// InstanceCommandRepository defines write operations for instances.
type InstanceCommandRepository interface {
	// Create persists a new instance within a transaction.
	Create(
		ctx context.Context,
		transaction shared.Transaction,
		instance *Instance,
	) (*Instance, error)

	// Update modifies an existing instance within a transaction.
	Update(
		ctx context.Context,
		transaction shared.Transaction,
		instance *Instance,
	) (*Instance, error)

	// Delete removes an instance by its ID within a transaction.
	Delete(
		ctx context.Context,
		transaction shared.Transaction,
		instanceID shared.InstanceID,
	) error
}

// InstanceQueryRepository defines read-only operations for instances.
type InstanceQueryRepository interface {
	// GetByID retrieves an instance by its ID.
	GetByID(
		ctx context.Context,
		instanceID shared.InstanceID,
	) (*Instance, error)

	// FetchNextBatchByListing returns instances for a given listing in batches.
	FetchNextBatchByListing(
		ctx context.Context,
		listingID shared.ListingID,
		request shared.BatchRequest,
	) (instances []*Instance, nextCursor shared.Cursor, err error)

	// FetchNextBatchByBillingAccount returns instances for a billing account in batches.
	FetchNextBatchByBillingAccount(
		ctx context.Context,
		billingAccountID shared.BillingAccountID,
		request shared.BatchRequest,
	) (instances []*Instance, nextCursor shared.Cursor, err error)

	// FetchNextBatchExpiredBefore returns instances expired before a cutoff time in batches.
	FetchNextBatchExpiredBefore(
		ctx context.Context,
		cutoffTime time.Time,
		request shared.BatchRequest,
	) (instances []*Instance, nextCursor shared.Cursor, err error)

	// FetchNextBatchByState returns instances filtered by instance and contract state.
	FetchNextBatchByState(
		ctx context.Context,
		instanceState InstanceState,
		contractState ContractState,
		request shared.BatchRequest,
	) (instances []*Instance, nextCursor shared.Cursor, err error)

	// FetchNextBatchPendingHardwareUpdate returns instances pending hardware updates.
	FetchNextBatchPendingHardwareUpdate(
		ctx context.Context,
		request shared.BatchRequest,
	) (instances []*Instance, nextCursor shared.Cursor, err error)
}
