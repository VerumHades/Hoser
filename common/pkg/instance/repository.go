package instance

// InstanceRepository defines persistence operations for instances.
type InstanceRepository interface {
	Save(instance *Instance) error
	GetByID(id string) (*Instance, error)
	Delete(id string) error
	ListByListing(listingID string) ([]*Instance, error)
	ListByBilling(billingID string) ([]*Instance, error)
}
