package billedcharge

import (
	"time"

	"common/pkg/currency"
	"common/pkg/hardware"
)

// BilledCharge represents a billed charge for a user and instance.
type BilledCharge struct {
	id                    string
	userID                string
	listingID             string
	instanceID            string
	pricePaid             currency.Money
	hardwareSpecification *hardware.HardwareSpecification
	endDate               time.Time
}

func (b *BilledCharge) ID() string                                    { return b.id }
func (b *BilledCharge) UserID() string                                { return b.userID }
func (b *BilledCharge) ListingID() string                             { return b.listingID }
func (b *BilledCharge) InstanceID() string                            { return b.instanceID }
func (b *BilledCharge) PricePaid() currency.Money                     { return b.pricePaid }
func (b *BilledCharge) HardwareSpec() *hardware.HardwareSpecification { return b.hardwareSpecification }
func (b *BilledCharge) EndDate() time.Time                            { return b.endDate }

func (b *BilledCharge) SetPricePaid(price currency.Money) { b.pricePaid = price }
func (b *BilledCharge) SetHardwareSpec(spec *hardware.HardwareSpecification) {
	b.hardwareSpecification = spec
}
func (b *BilledCharge) SetEndDate(end time.Time) { b.endDate = end }

func (b *BilledCharge) DaysRemaining() int {
	days := int(b.endDate.Sub(time.Now()).Hours() / 24)
	if days < 0 {
		return 0
	}
	return days
}
