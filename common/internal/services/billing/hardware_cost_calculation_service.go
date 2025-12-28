package application

import (
	"context"
	"time"

	"common/internal/domain/rates"
	"common/internal/shared"
)

// HardwareCostCalculationService calculates hardware costs over time using rates from a repository.
type HardwareCostCalculationService struct {
	rateRepository rates.HardwareCostQueryRepository
}

// NewHardwareCostCalculationService creates a new calculation service.
func NewHardwareCostCalculationService(
	rateRepository rates.HardwareCostQueryRepository,
) *HardwareCostCalculationService {
	return &HardwareCostCalculationService{
		rateRepository: rateRepository,
	}
}

// CalculateCost computes the total cost (in minor units) for a hardware specification over a duration.
func (s *HardwareCostCalculationService) CalculateCost(
	ctx context.Context,
	spec *shared.HardwareSpecification,
	duration time.Duration,
	at time.Time,
) (int64, error) {
	// Convert duration to hours, rounding up partial hours
	hours := int64(duration.Hours())
	if hours == 0 {
		hours = 1
	}

	totalCost := int64(0)

	// CPU
	cpuRate, err := s.rateRepository.GetActiveRate(ctx, rates.ResourceCPU, at)
	if err != nil {
		return 0, err
	}
	totalCost += cpuRate.CostInCents() * spec.CPUCount * hours

	// RAM
	ramRate, err := s.rateRepository.GetActiveRate(ctx, rates.ResourceRAM, at)
	if err != nil {
		return 0, err
	}
	totalCost += ramRate.CostInCents() * spec.RAMBytes * hours

	// Disk
	diskRate, err := s.rateRepository.GetActiveRate(ctx, rates.ResourceDisk, at)
	if err != nil {
		return 0, err
	}
	totalCost += diskRate.CostInCents() * spec.DiskBytes * hours

	return totalCost, nil
}
