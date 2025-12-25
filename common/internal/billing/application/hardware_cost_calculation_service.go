package application

import (
	"time"

	"common/internal/billing/domain"
	"common/internal/shared"
)

// HardwareCostCalculationService calculates hardware costs over time using rates from a repository.
type HardwareCostCalculationService struct {
	rateRepository    domain.HardwareCostRepository
	conversionService shared.CurrencyConversionService
}

// NewHardwareCostCalculationService creates a new calculation service.
func NewHardwareCostCalculationService(
	rateRepository domain.HardwareCostRepository,
	conversionService shared.CurrencyConversionService,
) *HardwareCostCalculationService {
	return &HardwareCostCalculationService{
		rateRepository:    rateRepository,
		conversionService: conversionService,
	}
}

// GetRateAt fetches the hardware cost rate at a given time.
func (s *HardwareCostCalculationService) GetRateAt(at time.Time) (*domain.HardwareCostRate, error) {
	return s.rateRepository.GetActiveRate(at)
}

// CalculateCost computes the total cost for a hardware specification over a given duration, optionally converting currency.
func (s *HardwareCostCalculationService) CalculateCost(
	spec *shared.HardwareSpecification,
	duration time.Duration,
	at time.Time,
	currencyCode string,
) (shared.Money, error) {

	rate, err := s.GetRateAt(at)
	if err != nil {
		return shared.Money{}, err
	}

	total := shared.Money{Amount: 0, CurrencyCode: rate.CPUCost().CurrencyCode}

	hours := duration.Hours()

	// CPU cost
	cpuCost := rate.CPUCost().Multiply(float64(spec.CPUCount) * hours)
	if currencyCode != "" && currencyCode != cpuCost.CurrencyCode {
		cpuCost, err = s.conversionService.Convert(cpuCost, currencyCode)
		if err != nil {
			return shared.Money{}, err
		}
	}
	total, err = total.Add(cpuCost, s.conversionService)
	if err != nil {
		return shared.Money{}, err
	}

	// RAM cost
	ramCost := rate.RAMCost().Multiply(float64(spec.RAMBytes) * hours)
	if currencyCode != "" && currencyCode != ramCost.CurrencyCode {
		ramCost, err = s.conversionService.Convert(ramCost, currencyCode)
		if err != nil {
			return shared.Money{}, err
		}
	}
	total, err = total.Add(ramCost, s.conversionService)
	if err != nil {
		return shared.Money{}, err
	}

	// Disk cost
	diskCost := rate.DiskCost().Multiply(float64(spec.DiskBytes) * hours)
	if currencyCode != "" && currencyCode != diskCost.CurrencyCode {
		diskCost, err = s.conversionService.Convert(diskCost, currencyCode)
		if err != nil {
			return shared.Money{}, err
		}
	}
	total, err = total.Add(diskCost, s.conversionService)
	if err != nil {
		return shared.Money{}, err
	}

	return total, nil
}
