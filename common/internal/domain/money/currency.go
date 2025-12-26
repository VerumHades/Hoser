package money

import (
	"common/internal/shared"
	"errors"
)

type CurrencyConversionService interface {
	Convert(money Money, targetCurrency Currency) (Money, error)
}

type Currency struct {
	code                   shared.CurrencyCode
	minorUnitDecimalPlaces int
}

var USD = Currency{code: "USD", minorUnitDecimalPlaces: 2}

// NewCurrency creates a validated Currency value object.
func NewCurrency(code shared.CurrencyCode, minorUnitDecimalPlaces int) (Currency, error) {
	if code == "" {
		return Currency{}, errors.New("currency code cannot be empty")
	}
	if minorUnitDecimalPlaces < 0 {
		return Currency{}, errors.New("minor unit decimal places cannot be negative")
	}

	return Currency{
		code:                   code,
		minorUnitDecimalPlaces: minorUnitDecimalPlaces,
	}, nil
}

func (currency Currency) Code() shared.CurrencyCode {
	return currency.code
}

func (currency Currency) MinorUnitDecimalPlaces() int {
	return currency.minorUnitDecimalPlaces
}

func (currency Currency) Equals(other Currency) bool {
	return currency.code == other.code
}
