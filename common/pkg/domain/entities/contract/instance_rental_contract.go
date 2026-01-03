package contract

import (
	"errors"
	"time"

	"common/pkg/shared"
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

	createdAt time.Time
}

// NewInstanceRentalContract creates a new contract period describing desired runtime state.
// newContract is a private initializer used by all constructors and mutators.
func newContract(
	id shared.InstanceRentalContractID,
	listingID shared.ListingID,
	ownerID shared.UserID,
	hardwareSpecification *shared.HardwareSpecification,
	periodStart, periodEnd time.Time,
	renewedFromContractID *shared.InstanceRentalContractID,
	canceledAt *time.Time,
	renewAutomatically bool,
	renewalDuration time.Duration,
	createdAt time.Time,
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
		createdAt:             createdAt,
	}

	if err := contract.Validate(); err != nil {
		return nil, err
	}

	return contract, nil
}

// NewInstanceRentalContract creates a new rental contract with a generated ID and current timestamp.
func NewInstanceRentalContract(
	listingID shared.ListingID,
	ownerID shared.UserID,
	hardwareSpecification *shared.HardwareSpecification,
	periodStart, periodEnd time.Time,
	renewAutomatically bool,
	renewalDuration time.Duration,
) (*InstanceRentalContract, error) {

	return newContract(
		shared.InstanceRentalContractID(shared.GenerateUUID()),
		listingID,
		ownerID,
		hardwareSpecification,
		periodStart,
		periodEnd,
		nil, // renewedFromContractID
		nil, // canceledAt
		renewAutomatically,
		renewalDuration,
		time.Now().UTC(),
	)
}

// NewInstanceRentalContractWithID creates a contract with a predefined ID and creation time.
func NewInstanceRentalContractWithID(
	id shared.InstanceRentalContractID,
	listingID shared.ListingID,
	ownerID shared.UserID,
	hardwareSpecification *shared.HardwareSpecification,
	periodStart, periodEnd time.Time,
	renewedFromContractID *shared.InstanceRentalContractID,
	canceledAt *time.Time,
	renewAutomatically bool,
	renewalDuration time.Duration,
	createdAt time.Time,
) (*InstanceRentalContract, error) {

	return newContract(
		id,
		listingID,
		ownerID,
		hardwareSpecification,
		periodStart,
		periodEnd,
		renewedFromContractID,
		canceledAt,
		renewAutomatically,
		renewalDuration,
		createdAt,
	)
}

// Renew derives a new contract from the current one, starting at the previous contract's end.
func (contract *InstanceRentalContract) Renew(now time.Time) (*InstanceRentalContract, error) {
	if contract == nil {
		return nil, errors.New("cannot renew a nil contract")
	}

	return newContract(
		shared.InstanceRentalContractID(shared.GenerateUUID()),
		contract.listingID,
		contract.ownerID,
		contract.hardwareSpecification,
		contract.periodEnd,
		contract.periodEnd.Add(contract.renewalDuration),
		&contract.id,
		nil, // canceledAt
		contract.renewAutomatically,
		contract.renewalDuration,
		time.Now().UTC(),
	)
}

// WithNewHardwareSpecification creates a new contract with updated hardware specification.
func (contract *InstanceRentalContract) WithNewHardwareSpecification(
	hardwareSpecification *shared.HardwareSpecification,
) (*InstanceRentalContract, error) {
	if contract == nil {
		return nil, errors.New("existing contract cannot be nil")
	}

	return newContract(
		shared.InstanceRentalContractID(shared.GenerateUUID()),
		contract.listingID,
		contract.ownerID,
		hardwareSpecification,
		contract.periodStart,
		contract.periodEnd,
		&contract.id,
		nil, // canceledAt
		contract.renewAutomatically,
		contract.renewalDuration,
		time.Now().UTC(),
	)
}
