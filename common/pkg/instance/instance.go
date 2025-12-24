package instance

import (
	"common/pkg/hardware"
	"common/pkg/util"
	"time"
)

// InstanceState defines the possible states of an instance.
type InstanceState int

const (
	Running InstanceState = iota
	Stopped

	BillingFailed
	ErrorState
)

type ContractState int

const (
	ContractInactive ContractState = iota
	ContractActive
)

/*
	ContractState InstanceState   Reaction

	Inactive      Running         The stop signal is given to the deployer
	Active        Stopped         The start signal is give to the deployer
	Active        BillingFailed   Retry payment, sets ContractState to Inactive on failiure
*/

type Instance struct {
	ID        string
	ListingID string
	BillingID string

	State         InstanceState
	ContractState ContractState

	HardwareSpecification *hardware.HardwareSpecification
	UpdateHardware        bool

	Expiry time.Time

	RenewAutomatically bool
	RenewalDuration    time.Duration
}

// InstanceRepository defines persistence operations for instances.
type InstanceRepository interface {
	Save(instance *Instance) error
	GetByID(id string) (*Instance, error)
	Delete(id string) error

	ListByListing(listingID string) ([]*Instance, error)
	ListByBilling(billingID string) ([]*Instance, error)

	ListExpired(cutoff time.Time) ([]*Instance, error)
	ListByState(state InstanceState, contractState ContractState) ([]*Instance, error)
	ListByPendingHardwareUpdate() ([]*Instance, error)
}

// InstanceService provides operations on instances.
type InstanceService struct {
	repo InstanceRepository
}

// NewInstanceService creates a new instance service.
func NewInstanceService(repo InstanceRepository) *InstanceService {
	return &InstanceService{repo: repo}
}

// CreateInstance creates a new instance from listing and billing IDs.
func (s *InstanceService) CreateInstance(listingID, billingID string, hardwareSpec *hardware.HardwareSpecification) (*Instance, error) {
	instance := &Instance{
		ID:                    util.GenerateUUID(),
		ListingID:             listingID,
		BillingID:             billingID,
		ContractState:         ContractActive,
		State:                 Stopped, // default to Running since Building isn't defined
		HardwareSpecification: hardwareSpec,
		RenewalDuration:       time.Hour,
	}
	if err := s.repo.Save(instance); err != nil {
		return nil, err
	}
	return instance, nil
}

// SetConcreteHardware sets the hardware and clears any desired spec.
func (s *InstanceService) SetHardware(instanceID string, hardwareSpec *hardware.HardwareSpecification) (*Instance, error) {
	instance, err := s.repo.GetByID(instanceID)
	if err != nil {
		return nil, err
	}

	instance.HardwareSpecification = hardwareSpec
	instance.UpdateHardware = true

	if err := s.repo.Save(instance); err != nil {
		return nil, err
	}
	return instance, nil
}

// SetExpiry updates the expiry date.
func (s *InstanceService) SetExpiry(instanceID string, expiry time.Time) (*Instance, error) {
	instance, err := s.repo.GetByID(instanceID)
	if err != nil {
		return nil, err
	}
	instance.Expiry = expiry
	if err := s.repo.Save(instance); err != nil {
		return nil, err
	}
	return instance, nil
}

// SetState updates the instance state.
func (s *InstanceService) SetState(instanceID string, state InstanceState) (*Instance, error) {
	instance, err := s.repo.GetByID(instanceID)
	if err != nil {
		return nil, err
	}
	instance.State = state
	if err := s.repo.Save(instance); err != nil {
		return nil, err
	}
	return instance, nil
}

// SetContractState updates the contract state of the instance.
func (s *InstanceService) SetContractState(instanceID string, contractState ContractState) (*Instance, error) {
	instance, err := s.repo.GetByID(instanceID)
	if err != nil {
		return nil, err
	}
	instance.ContractState = contractState
	if err := s.repo.Save(instance); err != nil {
		return nil, err
	}
	return instance, nil
}

// GetInstance retrieves an instance by ID.
func (s *InstanceService) GetInstance(instanceID string) (*Instance, error) {
	return s.repo.GetByID(instanceID)
}

// ListByListing retrieves all instances for a listing.
func (s *InstanceService) ListByListing(listingID string) ([]*Instance, error) {
	return s.repo.ListByListing(listingID)
}

// ListByBillingAccount retrieves all instances for a billing account.
func (s *InstanceService) ListByBillingAccount(billingID string) ([]*Instance, error) {
	return s.repo.ListByBilling(billingID)
}

// ListExpiredInstances retrieves expired instances.
func (s *InstanceService) ListExpiredInstances(cutoff time.Time) ([]*Instance, error) {
	return s.repo.ListExpired(cutoff)
}

// DeleteInstance deletes an instance.
func (s *InstanceService) DeleteInstance(instanceID string) error {
	return s.repo.Delete(instanceID)
}

// ListActiveNonRunning retrieves ContractActive instances that are not running.
func (s *InstanceService) ListActiveNonRunning() ([]*Instance, error) {
	return s.repo.ListByState(Stopped, ContractActive)
}

// ListInactiveRunning retrieves ContractInactive instances that are running.
func (s *InstanceService) ListInactiveRunning() ([]*Instance, error) {
	return s.repo.ListByState(Running, ContractInactive)
}

// ListInactiveRunning retrieves ContractInactive instances that are running.
func (s *InstanceService) ListByPendingHardwareUpdate() ([]*Instance, error) {
	return s.repo.ListByPendingHardwareUpdate()
}
