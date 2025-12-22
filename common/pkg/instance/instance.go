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
	if err := s.repo.Save(instance); err != nil {
		return nil, err
	}
	return instance, nil
}

// UpdateHardwareSpecification updates an instance's hardware spec.
func (s *InstanceService) UpdateHardwareSpecification(instanceID string, newHardwareSpec *hardware.HardwareSpecification) (*Instance, error) {
	instance, err := s.repo.GetByID(instanceID)
	if err != nil {
		return nil, err
	}
	instance.HardwareSpecification = newHardwareSpec
	if err := s.repo.Save(instance); err != nil {
		return nil, err
	}
	return instance, nil
}

// GetInstance retrieves an instance by ID.
func (s *InstanceService) GetInstance(instanceID string) (*Instance, error) {
	return s.repo.GetByID(instanceID)
}

// ListByBillingAccount lists all instances associated with a billing account.
func (s *InstanceService) ListByBillingAccount(billingID string) ([]*Instance, error) {
	return s.repo.ListByBilling(billingID)
}
