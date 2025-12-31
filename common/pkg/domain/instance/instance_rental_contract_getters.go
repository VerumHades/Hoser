package instance

import (
	"common/pkg/shared"
	"errors"
	"time"
)

// IsActiveAt reports whether the contract should be enforced at the given time.
func (contract *InstanceRentalContract) IsActiveAt(at time.Time) bool {
	if at.Before(contract.periodStart) || !at.Before(contract.periodEnd) {
		return false
	}
	if contract.canceledAt != nil && !at.Before(*contract.canceledAt) {
		return false
	}
	return true
}

// ShouldExistNow reports whether the reconciler should ensure runtime presence now.
func (contract *InstanceRentalContract) ShouldExistNow(now time.Time) bool {
	return contract.IsActiveAt(now)
}

// Validate ensures the contract describes a valid authorization window.
func (contract *InstanceRentalContract) Validate() error {
	if contract.id == "" {
		return errors.New("contract ID cannot be empty")
	}
	if contract.listingID == "" {
		return errors.New("listing ID cannot be empty")
	}
	if contract.ownerID == "" {
		return errors.New("owner ID cannot be empty")
	}
	if contract.hardwareSpecification == nil {
		return errors.New("hardware specification cannot be nil")
	}
	if contract.periodStart.IsZero() || contract.periodEnd.IsZero() {
		return errors.New("contract period cannot be zero")
	}
	if !contract.periodEnd.After(contract.periodStart) {
		return errors.New("contract period end must be after start")
	}
	if contract.renewAutomatically && contract.renewalDuration <= 0 {
		return errors.New("renewal duration must be positive when auto-renew is enabled")
	}
	return nil
}

// ID returns the unique identifier of the rental contract.
func (contract *InstanceRentalContract) ID() shared.InstanceRentalContractID {
	return contract.id
}

// ListingID returns the listing ID associated with this contract.
func (contract *InstanceRentalContract) ListingID() shared.ListingID {
	return contract.listingID
}

// OwnerID returns the owner ID of this contract.
func (contract *InstanceRentalContract) OwnerID() shared.UserID {
	return contract.ownerID
}

// HardwareSpecification returns the hardware specification for this contract.
func (contract *InstanceRentalContract) HardwareSpecification() *shared.HardwareSpecification {
	return contract.hardwareSpecification
}

// PeriodStart returns the start time of the rental period.
func (contract *InstanceRentalContract) PeriodStart() time.Time {
	return contract.periodStart
}

// PeriodEnd returns the end time of the rental period.
func (contract *InstanceRentalContract) PeriodEnd() time.Time {
	return contract.periodEnd
}

// RenewedFromContractID returns the ID of the contract this one was renewed from, if any.
func (contract *InstanceRentalContract) RenewedFromContractID() *shared.InstanceRentalContractID {
	return contract.renewedFromContractID
}

// CanceledAt returns the cancellation timestamp, if the contract was canceled.
func (contract *InstanceRentalContract) CanceledAt() *time.Time {
	return contract.canceledAt
}

// CreatedAt returns the timestamp when the contract was created.
func (contract *InstanceRentalContract) CreatedAt() time.Time {
	return contract.createdAt
}

// RenewAutomatically indicates whether the contract is set to auto-renew.
func (contract *InstanceRentalContract) RenewAutomatically() bool {
	return contract.renewAutomatically
}

// RenewalDuration returns the duration for automatic renewal, if enabled.
func (contract *InstanceRentalContract) RenewalDuration() time.Duration {
	return contract.renewalDuration
}
