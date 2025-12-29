package instance

import (
	"errors"
	"time"

	"common/internal/shared"
)

type InstanceRentalContract struct {
	id shared.InstanceRentalContractID

	listingID shared.ListingID
	ownerID   shared.UserID

	hardwareSpecification *shared.HardwareSpecification

	periodStart           time.Time
	periodEnd             time.Time
	renewedFromContractID *shared.InstanceRentalContractID

	canceledAt *time.Time

	renewAutomatically bool
	renewalDuration    time.Duration
}

// NewInstanceRentalContract creates a new contract period describing desired runtime state.
func NewInstanceRentalContract(
	listingID shared.ListingID,
	ownerID shared.UserID,
	hardwareSpecification *shared.HardwareSpecification,
	periodStart time.Time,
	periodEnd time.Time,
	renewAutomatically bool,
	renewalDuration time.Duration,
) (*InstanceRentalContract, error) {

	contract := &InstanceRentalContract{
		id:                    shared.InstanceRentalContractID(shared.GenerateUUID()),
		listingID:             listingID,
		ownerID:               ownerID,
		hardwareSpecification: hardwareSpecification,
		periodStart:           periodStart,
		periodEnd:             periodEnd,
		renewAutomatically:    renewAutomatically,
		renewalDuration:       renewalDuration,
	}

	if err := contract.Validate(); err != nil {
		return nil, err
	}

	return contract, nil
}

func NewInstanceRentalContractWithID(
	id shared.InstanceRentalContractID,
	listingID shared.ListingID,
	ownerID shared.UserID,
	hardwareSpecification *shared.HardwareSpecification,
	periodStart time.Time,
	periodEnd time.Time,
	renewedFromContractID *shared.InstanceRentalContractID,
	canceledAt *time.Time,
	renewAutomatically bool,
	renewalDuration time.Duration,
) (*InstanceRentalContract, error) {
	contract := &InstanceRentalContract{
		id:                    id,
		listingID:             listingID,
		ownerID:               ownerID,
		hardwareSpecification: hardwareSpecification,
		periodStart:           periodStart,
		periodEnd:             periodEnd,
		renewedFromContractID: renewedFromContractID,
		canceledAt:            canceledAt,
		renewAutomatically:    renewAutomatically,
		renewalDuration:       renewalDuration,
	}

	if err := contract.Validate(); err != nil {
		return nil, err
	}

	return contract, nil
}

// Renew derives a new contract from this one if renewal is allowed.
func (contract *InstanceRentalContract) Renew(now time.Time) (*InstanceRentalContract, error) {
	newContract := &InstanceRentalContract{
		id:                    shared.InstanceRentalContractID(shared.GenerateUUID()),
		listingID:             contract.listingID,
		ownerID:               contract.ownerID,
		hardwareSpecification: contract.hardwareSpecification,
		periodStart:           contract.periodEnd,
		periodEnd:             contract.periodEnd.Add(contract.renewalDuration),
		renewedFromContractID: &contract.id,
		renewAutomatically:    contract.renewAutomatically,
		renewalDuration:       contract.renewalDuration,
	}

	if err := newContract.Validate(); err != nil {
		return nil, err
	}

	return newContract, nil
}

func (existingContract *InstanceRentalContract) WithNewHardwareSpecification(
	hardwareSpecification *shared.HardwareSpecification,
) (*InstanceRentalContract, error) {
	if existingContract == nil {
		return nil, errors.New("existing contract cannot be nil")
	}

	newContract := &InstanceRentalContract{
		id:                    shared.InstanceRentalContractID(shared.GenerateUUID()),
		listingID:             existingContract.listingID,
		ownerID:               existingContract.ownerID,
		hardwareSpecification: hardwareSpecification,
		periodStart:           existingContract.periodStart,
		periodEnd:             existingContract.periodEnd,
		renewedFromContractID: &existingContract.id,
		renewAutomatically:    existingContract.renewAutomatically,
		renewalDuration:       existingContract.renewalDuration,
	}

	if err := newContract.Validate(); err != nil {
		return nil, err
	}

	return newContract, nil
}

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
