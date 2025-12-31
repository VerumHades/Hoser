package rates

import (
	"context"
	"time"

	"common/pkg/shared"
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

// / HardwareCostRateCursor represents a stable pagination position
// / for iterating over hardware cost rates.
type HardwareCostRateCursor struct {
	LastEffectiveDate time.Time
}

// HardwareCostQueryRepository defines read-only operations for hardware cost rates.
type HardwareCostQueryRepository interface {
	GetActiveRate(
		ctx context.Context,
		resourceType HardwareResourceType,
		at time.Time,
	) (*HardwareCostRate, error)

	FetchNextBatchOrderedByEffectiveDate(
		ctx context.Context,
		request shared.BatchRequest[HardwareCostRateCursor],
	) (rates []*HardwareCostRate, nextCursor HardwareCostRateCursor, err error)
}
