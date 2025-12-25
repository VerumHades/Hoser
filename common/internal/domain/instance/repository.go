package instance

import (
	"common/internal/shared"
	"time"
)

// InstanceRepository defines persistence operations for instances.
type InstanceRepository interface {
	Save(instance *Instance) error
	GetByID(id shared.InstanceID) (*Instance, error)
	Delete(id shared.InstanceID) error

	ListByListing(listingID shared.ListingID) ([]*Instance, error)
	ListByBilling(billingID shared.BillingAccountID) ([]*Instance, error)

	ListExpired(cutoff time.Time) ([]*Instance, error)
	ListByState(state InstanceState, contractState ContractState) ([]*Instance, error)
	ListByPendingHardwareUpdate() ([]*Instance, error)
}
