package inmemconversion

import (
	"common/pkg/money"
	"fmt"
)

// DummyCurrencyConversionService is a simple implementation for testing.
type DummyCurrencyConversionService struct {
	rates map[string]float64 // key: "USD:EUR", value: rate
}

// NewDummyCurrencyConversionService creates a new dummy conversion service with preset rates.
func NewDummyCurrencyConversionService() *DummyCurrencyConversionService {
	return &DummyCurrencyConversionService{
		rates: map[string]float64{
			"USD:EUR": 0.9,
			"EUR:USD": 1.1,
			"USD:USD": 1.0,
			"EUR:EUR": 1.0,
		},
	}
}

// Convert converts the given Money to the target currency using fixed rates.
func (d *DummyCurrencyConversionService) Convert(m money.Money, targetCurrency string) (money.Money, error) {
	if m.CurrencyCode == targetCurrency {
		return m, nil
	}

	key := fmt.Sprintf("%s:%s", m.CurrencyCode, targetCurrency)
	rate, ok := d.rates[key]
	if !ok {
		return money.Money{}, fmt.Errorf("no conversion rate from %s to %s", m.CurrencyCode, targetCurrency)
	}

	return money.Money{
		Amount:       m.Amount * rate,
		CurrencyCode: targetCurrency,
	}, nil
}
