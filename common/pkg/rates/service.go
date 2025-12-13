package rates

import (
	"time"
)

// HardwareCostService provides access to hardware cost rates.
type HardwareCostService struct {
	repo HardwareCostRepository
}

func NewHardwareCostService(repo HardwareCostRepository) *HardwareCostService {
	return &HardwareCostService{repo: repo}
}

// GetRateView returns the hardware cost rate active at the specified time.
func (s *HardwareCostService) GetRateView(at time.Time) (*HardwareCostRateView, error) {
	rate, err := s.repo.GetActiveRate(at)
	if err != nil {
		return nil, err
	}

	return &HardwareCostRateView{
		ID:        rate.ID,
		CPUCost:   rate.CPUCost,
		RAMCost:   rate.RAMCost,
		DiskCost:  rate.DiskCost,
		ValidFrom: rate.ValidFrom,
		ValidTo:   rate.ValidTo,
	}, nil
}
