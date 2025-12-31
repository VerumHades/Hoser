package instance

import (
	"errors"
	"time"
)

// EnableRenewal sets the contract to renew automatically with the specified duration.
func (contract *InstanceRentalContract) EnableRenewal(renewalDuration time.Duration) error {
	if renewalDuration <= 0 {
		return errors.New("renewal duration must be positive when enabling auto-renew")
	}

	contract.renewAutomatically = true
	contract.renewalDuration = renewalDuration
	return nil
}

// DisableRenewal disables automatic renewal for the contract.
func (contract *InstanceRentalContract) DisableRenewal() {
	contract.renewAutomatically = false
	contract.renewalDuration = 0
}

// Cancel marks the contract as no longer authoritative from the given time.
func (contract *InstanceRentalContract) Cancel(at time.Time) error {
	if contract.canceledAt != nil {
		return errors.New("contract already canceled")
	}

	if at.Before(contract.periodStart) {
		return errors.New("cancellation cannot precede contract start")
	}

	contract.canceledAt = &at
	return nil
}
