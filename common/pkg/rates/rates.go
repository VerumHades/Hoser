package rates

import (
	"common/pkg/money"
	"time"

	"github.com/google/uuid"
)

// HardwareCostRate represents the cost of CPU, RAM, and Disk over a period.
type HardwareCostRate struct {
	ID        string
	CPUCost   money.Money
	RAMCost   money.Money
	DiskCost  money.Money
	ValidFrom time.Time
	ValidTo   *time.Time // nil = currently valid indefinitely
}

// HardwareCostRepository defines persistence operations for hardware cost rates.
type HardwareCostRepository interface {
	Save(rate *HardwareCostRate) error
	GetActiveRate(at time.Time) (*HardwareCostRate, error)
	ListAll() ([]*HardwareCostRate, error)
}

// HardwareCostService provides access to hardware cost rates.
type HardwareCostService struct {
	repo HardwareCostRepository
}

// NewHardwareCostService creates a new service instance.
func NewHardwareCostService(repo HardwareCostRepository) *HardwareCostService {
	return &HardwareCostService{repo: repo}
}

// GetRate returns the hardware cost rate active at the specified time.
func (s *HardwareCostService) GetRate(at time.Time) (*HardwareCostRate, error) {
	return s.repo.GetActiveRate(at)
}

// AddRate creates and saves a new hardware cost rate.
func (s *HardwareCostService) AddRate(cpuCost, ramCost, diskCost money.Money, validFrom time.Time, validTo *time.Time) (*HardwareCostRate, error) {
	rate := &HardwareCostRate{
		ID:        uuid.NewString(),
		CPUCost:   cpuCost,
		RAMCost:   ramCost,
		DiskCost:  diskCost,
		ValidFrom: validFrom,
		ValidTo:   validTo,
	}
	if err := s.repo.Save(rate); err != nil {
		return nil, err
	}
	return rate, nil
}
