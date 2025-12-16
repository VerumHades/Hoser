package rates

import (
	"common/pkg/money"
	"time"
)

type HardwareCostRate struct {
	ID        string
	CPUCost   money.Money
	RAMCost   money.Money
	DiskCost  money.Money
	ValidFrom time.Time
	ValidTo   *time.Time // nil = currently valid indefinitely
}

type HardwareCostRepository interface {
	Save(rate *HardwareCostRate) error
	GetActiveRate(at time.Time) (*HardwareCostRate, error)
	ListAll() ([]*HardwareCostRate, error)
}

// HardwareCostService provides access to hardware cost rates.
type HardwareCostService struct {
	repo HardwareCostRepository
}

func NewHardwareCostService(repo HardwareCostRepository) *HardwareCostService {
	return &HardwareCostService{repo: repo}
}

// GetRateView returns the hardware cost rate active at the specified time.
func (s *HardwareCostService) GetRate(at time.Time) (*HardwareCostRate, error) {
	return s.repo.GetActiveRate(at)
}
