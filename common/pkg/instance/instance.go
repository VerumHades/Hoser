package instance

import (
	"common/pkg/hardware"
	"common/pkg/util"
)

type InstanceState int

const (
	Building InstanceState = iota
	Running
	Stopped
	NotBilled
)

// Instance represents a hardware rental instance derived from a listing.
type Instance struct {
	ID                    string
	ListingID             string
	BillingID             string
	State                 InstanceState
	HardwareSpecification *hardware.HardwareSpecification
}

// InstanceRepository defines persistence operations for instances.
type InstanceRepository interface {
	Save(instance *Instance) error
	GetByID(id string) (*Instance, error)
	Delete(id string) error
	ListByListing(listingID string) ([]*Instance, error)
	ListByBilling(billingID string) ([]*Instance, error)
}

// InstanceService orchestrates business logic for instances.
type InstanceService struct {
	repo InstanceRepository
}

// NewInstanceService creates a new service instance.
func NewInstanceService(repo InstanceRepository) *InstanceService {
	return &InstanceService{repo: repo}
}

// CreateInstance creates a new instance from a listing and billing ID.
func (s *InstanceService) CreateInstance(listingID, billingID string, hardwareSpec *hardware.HardwareSpecification) (*Instance, error) {
	instance := &Instance{
		ID:                    util.GenerateUUID(),
		ListingID:             listingID,
		BillingID:             billingID,
		State:                 Building,
		HardwareSpecification: hardwareSpec,
	}
	err := s.repo.Save(instance)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

// UpdateState updates the state of an instance.
func (s *InstanceService) UpdateState(instanceID string, state InstanceState) error {
	instance, err := s.repo.GetByID(instanceID)
	if err != nil {
		return err
	}
	instance.State = state
	return s.repo.Save(instance)
}

// GetInstanceView retrieves an instance as a read-only view.
func (s *InstanceService) GetInstance(instanceID string) (*Instance, error) {
	return s.repo.GetByID(instanceID)
}
