package application

import (
	"time"

	"common/internal/domain/billing"
	"common/internal/domain/money"

	"common/internal/shared"
)

// HardwareCostCalculationService calculates hardware costs over time using rates from a repository.
type HardwareCostCalculationService struct {
	rateRepository     billing.HardwareCostRepository
	currencyRepository money.CurrencyRepository
	conversionService  money.CurrencyConversionService
}

// NewHardwareCostCalculationService creates a new calculation service.
func NewHardwareCostCalculationService(
	rateRepository billing.HardwareCostRepository,
	conversionService money.CurrencyConversionService,
) *HardwareCostCalculationService {
	return &HardwareCostCalculationService{
		rateRepository:    rateRepository,
		conversionService: conversionService,
	}
}

// GetRateAt fetches the hardware cost rate at a given time.
func (s *HardwareCostCalculationService) GetRateAt(at time.Time) (*billing.HardwareCostRate, error) {
	return s.rateRepository.GetActiveRate(at)
}

// CalculateCost computes the total cost for a hardware specification over a given duration, optionally converting currency.
func (s *HardwareCostCalculationService) CalculateCost(
	spec *shared.HardwareSpecification,
	duration time.Duration,
	at time.Time,
	targetCurrencyCode shared.CurrencyCode,
) (money.Money, error) {

	rate, err := s.GetRateAt(at)
	if err != nil {
		return money.Money{}, err
	}

	targetCurrency, err := s.currencyRepository.FindByCode(targetCurrencyCode)
	if err != nil {
		return money.Money{}, err
	}
	totalCost, err := money.NewMoneyFromMinorUnits(0, targetCurrency)
	if err != nil {
		return money.Money{}, err
	}

	hours := duration.Hours()

	type resource struct {
		quantity int64
		baseCost money.Money
	}

	resources := []resource{
		{quantity: spec.CPUCount, baseCost: rate.CPUCost()},
		{quantity: spec.RAMBytes, baseCost: rate.RAMCost()},
		{quantity: spec.DiskBytes, baseCost: rate.DiskCost()},
	}

	for _, res := range resources {
		cost, err := res.baseCost.Multiply(res.quantity * int64(hours))
		if err != nil {
			return money.Money{}, err
		}

		totalCost, err = cost.Add(totalCost, s.conversionService)
		if err != nil {
			return money.Money{}, err
		}
	}

	return totalCost, nil
}
