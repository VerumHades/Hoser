package billedcharge

import "time"

// BilledChargeRepository defines persistence operations.
type BilledChargeRepository interface {
	Save(charge *BilledCharge) error
	GetByID(id string) (*BilledCharge, error)
	ListByUser(userID string) ([]*BilledCharge, error)
	ListByListing(listingID string) ([]*BilledCharge, error)
	ListByInstance(instanceID string) ([]*BilledCharge, error)
	ListActiveByUser(userID string, now time.Time) ([]*BilledCharge, error)
	Delete(id string) error
}
