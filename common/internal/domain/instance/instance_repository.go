package instance

import (
	"common/internal/shared"
	"time"
)

// InstanceRepository defines persistence operations for instances.
type InstanceRepository interface {
	Save(instance *Instance) (*Instance, error)

	GetByID(instanceID shared.InstanceID) (*Instance, error)
	Delete(instanceID shared.InstanceID) error

	FetchNextBatchByListing(
		listingID shared.ListingID,
		lastSeenInstanceID shared.InstanceID,
		maximumBatchSize int,
	) ([]*Instance, error)

	FetchNextBatchByBillingAccount(
		billingAccountID shared.BillingAccountID,
		lastSeenInstanceID shared.InstanceID,
		maximumBatchSize int,
	) ([]*Instance, error)

	FetchNextBatchExpiredBefore(
		cutoffTime time.Time,
		lastSeenInstanceID shared.InstanceID,
		maximumBatchSize int,
	) ([]*Instance, error)

	FetchNextBatchByState(
		instanceState InstanceState,
		contractState ContractState,
		lastSeenInstanceID shared.InstanceID,
		maximumBatchSize int,
	) ([]*Instance, error)

	FetchNextBatchPendingHardwareUpdate(
		lastSeenInstanceID shared.InstanceID,
		maximumBatchSize int,
	) ([]*Instance, error)
}
