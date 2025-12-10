package billedcharge

import (
	"common/pkg/currency"
	"common/pkg/hardware"
	"time"
)

// BilledChargeView is a read-only representation of a charge.
type BilledChargeView struct {
	ID                    string
	UserID                string
	ListingID             string
	InstanceID            string
	PricePaid             currency.Money
	HardwareSpecification *hardware.HardwareSpecification
	EndDate               time.Time
	DaysRemaining         int
}

// ToView converts a BilledCharge model to a read-only view.
func (b *BilledCharge) ToView() *BilledChargeView {
	return &BilledChargeView{
		ID:                    b.id,
		UserID:                b.userID,
		ListingID:             b.listingID,
		InstanceID:            b.instanceID,
		PricePaid:             b.pricePaid,
		HardwareSpecification: b.hardwareSpecification,
		EndDate:               b.endDate,
		DaysRemaining:         b.DaysRemaining(),
	}
}
