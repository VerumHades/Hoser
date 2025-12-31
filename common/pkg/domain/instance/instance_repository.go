package instance

import (
	"common/pkg/shared"
	"context"
	"time"
)

// InstanceRentalContractCommandRepository defines write operations for rental contracts.
type InstanceRentalContractCommandRepository interface {
	Create(
		ctx context.Context,
		transaction shared.Transaction,
		contract *InstanceRentalContract,
	) error

	Update(
		ctx context.Context,
		transaction shared.Transaction,
		contract *InstanceRentalContract,
	) error

	Delete(
		ctx context.Context,
		transaction shared.Transaction,
		contractID shared.InstanceRentalContractID,
	) error
}

// / InstanceRentalContractCursor represents a stable pagination position
// / for iterating over instance rental contracts.
type InstanceRentalContractCursor struct {
	LastCreatedAt time.Time
}

// InstanceRentalContractQueryRepository defines read-only operations for rental contracts.
type InstanceRentalContractQueryRepository interface {
	GetByID(
		ctx context.Context,
		contractID shared.InstanceRentalContractID,
	) (*InstanceRentalContract, error)

	Exists(
		ctx context.Context,
		contractID shared.InstanceRentalContractID,
	) (bool, error)

	FetchNextBatchByListing(
		ctx context.Context,
		listingID shared.ListingID,
		request shared.BatchRequest[InstanceRentalContractCursor],
	) (contracts []*InstanceRentalContract, nextCursor InstanceRentalContractCursor, err error)

	FetchNextBatchByOwner(
		ctx context.Context,
		ownerID shared.UserID,
		request shared.BatchRequest[InstanceRentalContractCursor],
	) (contracts []*InstanceRentalContract, nextCursor InstanceRentalContractCursor, err error)

	FetchNextBatchActiveAt(
		ctx context.Context,
		at time.Time,
		request shared.BatchRequest[InstanceRentalContractCursor],
	) (contracts []*InstanceRentalContract, nextCursor InstanceRentalContractCursor, err error)

	FetchNextBatchExpiredBefore(
		ctx context.Context,
		cutoffTime time.Time,
		request shared.BatchRequest[InstanceRentalContractCursor],
	) (contracts []*InstanceRentalContract, nextCursor InstanceRentalContractCursor, err error)

	FetchNextBatchPendingRenewal(
		ctx context.Context,
		cutoffTime time.Time,
		request shared.BatchRequest[InstanceRentalContractCursor],
	) (contracts []*InstanceRentalContract, nextCursor InstanceRentalContractCursor, err error)
}
