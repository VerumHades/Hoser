package domain

import (
	"errors"
	"time"

	"common/internal/shared"
)

// InstanceState defines the possible states of an instance.
type InstanceState int

const (
	Running InstanceState = iota
	Stopped
	BillingFailed
	ErrorState
)

// ContractState defines the possible states of the instance's contract.
type ContractState int

const (
	ContractInactive ContractState = iota
	ContractActive
)

// Instance represents a deployed resource linked to a listing and billing account.
type Instance struct {
	id        shared.InstanceID
	listingID shared.ListingID
	billingID shared.BillingAccountID

	state         InstanceState
	contractState ContractState

	hardwareSpecification *shared.HardwareSpecification
	updateHardware        bool

	expiry time.Time

	renewAutomatically bool
	renewalDuration    time.Duration
}

// NewInstance creates a new Instance with a generated ID.
func NewInstance(listingID shared.ListingID, billingID shared.BillingAccountID, hardwareSpec *shared.HardwareSpecification, expiry time.Time) (*Instance, error) {
	id := shared.InstanceID(shared.GenerateUUID())
	return NewInstanceWithID(id, listingID, billingID, hardwareSpec, expiry)
}

// NewInstanceWithID creates an Instance with an existing ID (for repository hydration).
func NewInstanceWithID(id shared.InstanceID, listingID shared.ListingID, billingID shared.BillingAccountID, hardwareSpec *shared.HardwareSpecification, expiry time.Time) (*Instance, error) {
	instance := &Instance{
		id:                    id,
		listingID:             listingID,
		billingID:             billingID,
		state:                 Stopped,
		contractState:         ContractInactive,
		hardwareSpecification: hardwareSpec,
		expiry:                expiry,
		renewAutomatically:    false,
		renewalDuration:       0,
	}
	if err := instance.Validate(); err != nil {
		return nil, err
	}
	return instance, nil
}

// ID returns the unique identifier.
func (i *Instance) ID() shared.InstanceID {
	return i.id
}

// shared.ListingID returns the associated listing ID.
func (i *Instance) ListingID() shared.ListingID {
	return i.listingID
}

// shared.BillingAccountID returns the associated billing ID.
func (i *Instance) BillingAccountID() shared.BillingAccountID {
	return i.billingID
}

// State returns the current instance state.
func (i *Instance) State() InstanceState {
	return i.state
}

// ContractState returns the current contract state.
func (i *Instance) ContractState() ContractState {
	return i.contractState
}

// Validate ensures the instance is in a valid state.
func (i *Instance) Validate() error {
	if i.id == "" {
		return errors.New("instance ID cannot be empty")
	}
	if i.listingID == "" {
		return errors.New("listing ID cannot be empty")
	}
	if i.billingID == "" {
		return errors.New("billing ID cannot be empty")
	}
	if i.hardwareSpecification == nil {
		return errors.New("shared specification cannot be nil")
	}
	if i.expiry.IsZero() {
		return errors.New("expiry cannot be zero")
	}
	return nil
}

// Start activates the instance and updates state.
func (i *Instance) Start() error {
	if i.state == Running {
		return errors.New("instance already running")
	}
	i.state = Running
	if i.contractState == ContractInactive {
		i.contractState = ContractActive
	}
	return nil
}

// Stop deactivates the instance and updates state.
func (i *Instance) Stop() error {
	if i.state == Stopped {
		return errors.New("instance already stopped")
	}
	i.state = Stopped
	return nil
}

// MarkBillingFailed sets the instance state to BillingFailed and optionally adjusts contract state.
func (i *Instance) MarkBillingFailed() {
	i.state = BillingFailed
	i.contractState = ContractInactive
}

// Renew attempts to renew the contract if auto-renew is enabled.
func (i *Instance) Renew() error {
	if !i.renewAutomatically {
		return errors.New("auto-renew not enabled")
	}
	i.expiry = i.expiry.Add(i.renewalDuration)
	i.contractState = ContractActive
	return nil
}

// EnableAutoRenew sets auto-renew and duration.
func (i *Instance) EnableAutoRenew(duration time.Duration) {
	i.renewAutomatically = true
	i.renewalDuration = duration
}

// DisableAutoRenew disables automatic renewal.
func (i *Instance) DisableAutoRenew() {
	i.renewAutomatically = false
	i.renewalDuration = 0
}
