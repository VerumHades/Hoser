package app

import (
	"common/pkg/hardware"
	"common/pkg/money"
	"common/pkg/rates"
	"time"
)

// HardwareCostCalculationService calculates hardware costs with optional currency conversions.
type HardwareCostCalculationService struct {
	costService       *rates.HardwareCostService
	conversionService money.CurrencyConversionService
}

// NewHardwareCostCalculationService creates a new calculation service.
func NewHardwareCostCalculationService(costService *rates.HardwareCostService, conversionService money.CurrencyConversionService) *HardwareCostCalculationService {
	return &HardwareCostCalculationService{
		costService:       costService,
		conversionService: conversionService,
	}
}

// GetRateView returns the hardware cost rate at a given time, optionally converted to another currency.
func (s *HardwareCostCalculationService) GetRateView(at time.Time, currencyCode string) (*rates.HardwareCostRate, error) {
	view, err := s.costService.GetRate(at)
	if err != nil {
		return nil, err
	}

	// Convert rates if a different currency is requested
	if currencyCode != "" && currencyCode != view.CPUCost.CurrencyCode {
		view.CPUCost, err = s.conversionService.Convert(view.CPUCost, currencyCode)
		if err != nil {
			return nil, err
		}
		view.RAMCost, err = s.conversionService.Convert(view.RAMCost, currencyCode)
		if err != nil {
			return nil, err
		}
		view.DiskCost, err = s.conversionService.Convert(view.DiskCost, currencyCode)
		if err != nil {
			return nil, err
		}
	}

	return view, nil
}

// CalculateCost calculates the total cost of a hardware specification over a given duration, optionally converting currency.
func (s *HardwareCostCalculationService) CalculateCost(
	spec *hardware.HardwareSpecification,
	duration time.Duration,
	at time.Time,
	currencyCode string,
) (money.Money, error) {

	view, err := s.GetRateView(at, currencyCode)
	if err != nil {
		return money.Money{}, err
	}

	total := money.Money{Amount: 0, CurrencyCode: currencyCode}

	hours := duration.Hours()

	cpuCost := view.CPUCost.Multiply(float64(spec.CPUCount) * hours)
	total, err = total.Add(cpuCost, s.conversionService)
	if err != nil {
		return money.Money{}, err
	}

	ramCost := view.RAMCost.Multiply(float64(spec.RAMBytes) * hours)
	total, err = total.Add(ramCost, s.conversionService)
	if err != nil {
		return money.Money{}, err
	}

	diskCost := view.DiskCost.Multiply(float64(spec.DiskBytes) * hours)
	total, err = total.Add(diskCost, s.conversionService)
	if err != nil {
		return money.Money{}, err
	}

	return total, nil
}
