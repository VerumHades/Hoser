package instance

import (
	"common/pkg/hardware"
	"common/pkg/util"
)

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
		id:                    util.GenerateUUID(),
		listingID:             listingID,
		billingID:             billingID,
		state:                 Building,
		hardwareSpecification: hardwareSpec,
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
	instance.state = state
	return s.repo.Save(instance)
}

// GetInstanceView retrieves an instance as a read-only view.
func (s *InstanceService) GetInstanceView(instanceID string) (*InstanceView, error) {
	instance, err := s.repo.GetByID(instanceID)
	if err != nil {
		return nil, err
	}
	return instance.ToView(), nil
}
