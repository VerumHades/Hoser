package billedcharge

import (
	"time"

	"common/pkg/currency"
	"common/pkg/hardware"
	"common/pkg/util"
)

// BilledChargeService handles business logic for billed charges.
type BilledChargeService struct {
	repo BilledChargeRepository
}

// NewBilledChargeService creates a new service instance.
func NewBilledChargeService(repo BilledChargeRepository) *BilledChargeService {
	return &BilledChargeService{repo: repo}
}

// CreateBilledCharge creates and saves a new billed charge.
func (s *BilledChargeService) CreateBilledCharge(userID, listingID, instanceID string, price currency.Money, spec *hardware.HardwareSpecification, endDate time.Time) (*BilledCharge, error) {
	charge := &BilledCharge{
		id:                    util.GenerateUUID(),
		userID:                userID,
		listingID:             listingID,
		instanceID:            instanceID,
		pricePaid:             price,
		hardwareSpecification: spec,
		endDate:               endDate,
	}
	err := s.repo.Save(charge)
	if err != nil {
		return nil, err
	}
	return charge, nil
}

// GetChargeView retrieves a read-only view of a billed charge.
func (s *BilledChargeService) GetChargeView(id string) (*BilledChargeView, error) {
	charge, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return charge.ToView(), nil
}

// ListActiveChargesByUser returns all currently active charges for a user.
func (s *BilledChargeService) ListActiveChargesByUser(userID string) ([]*BilledChargeView, error) {
	charges, err := s.repo.ListActiveByUser(userID, time.Now())
	if err != nil {
		return nil, err
	}

	views := make([]*BilledChargeView, len(charges))
	for i, c := range charges {
		views[i] = c.ToView()
	}
	return views, nil
}
