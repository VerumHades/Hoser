package instance

import "common/pkg/hardware"

// InstanceView is a read-only representation of an instance.
type InstanceView struct {
	ID                    string
	ListingID             string
	BillingID             string
	State                 InstanceState
	HardwareSpecification *hardware.HardwareSpecification
}

// ToView converts an Instance model to a view.
func (i *Instance) ToView() *InstanceView {
	return &InstanceView{
		ID:                    i.id,
		ListingID:             i.listingID,
		BillingID:             i.billingID,
		State:                 i.state,
		HardwareSpecification: i.hardwareSpecification,
	}
}
