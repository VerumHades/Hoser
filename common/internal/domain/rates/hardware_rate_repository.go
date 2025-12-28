package rates

import (
	"context"
	"time"

	"common/internal/shared"
)

// HardwareCostCommandRepository defines write operations for hardware cost rates.
type HardwareCostCommandRepository interface {
	// Create persists a new hardware cost rate within a transaction.
	Create(
		ctx context.Context,
		transaction shared.Transaction,
		rate *HardwareCostRate,
	) (*HardwareCostRate, error)

	// Update modifies an existing hardware cost rate within a transaction.
	// Typically used to correct or replace a rate; rates are otherwise immutable.
	Update(
		ctx context.Context,
		transaction shared.Transaction,
		rate *HardwareCostRate,
	) (*HardwareCostRate, error)
}

// HardwareCostQueryRepository defines read-only operations for hardware cost rates.
type HardwareCostQueryRepository interface {
	// GetActiveRate returns the currently active hardware cost rate for a given resource.
	GetActiveRate(
		ctx context.Context,
		resourceType HardwareResourceType,
		at time.Time,
	) (*HardwareCostRate, error)

	// FetchNextBatchOrderedByEffectiveDate returns hardware cost rates in order of effective date.
	// Useful for audits or migrations.
	FetchNextBatchOrderedByEffectiveDate(
		ctx context.Context,
		request shared.BatchRequest,
	) (rates []*HardwareCostRate, nextCursor shared.Cursor, err error)
}
