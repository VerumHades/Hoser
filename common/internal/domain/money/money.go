package money

import (
	"errors"
	"fmt"
)

type Money struct {
	amountInMinorUnits int64
	currency           Currency
}

// NewMoneyFromMinorUnits creates a Money instance using minor units (e.g. cents).
func NewMoneyFromMinorUnits(
	amountInMinorUnits int64,
	currency Currency,
) (Money, error) {
	if amountInMinorUnits < 0 {
		return Money{}, errors.New("amount cannot be negative")
	}

	return Money{
		amountInMinorUnits: amountInMinorUnits,
		currency:           currency,
	}, nil
}

func (money Money) AmountInMinorUnits() int64 {
	return money.amountInMinorUnits
}

func (money Money) Currency() Currency {
	return money.currency
}

func (money Money) IsZero() bool {
	return money.amountInMinorUnits == 0
}

func (money Money) operationWithNormalizedCurrency(
	other Money,
	conversionService CurrencyConversionService,
	operation func(a int64, b int64) int64,
) (Money, error) {
	normalizedOther, err := money.normalizeCurrency(other, conversionService)
	if err != nil {
		return Money{}, err
	}

	resultAmount := operation(money.amountInMinorUnits, normalizedOther.amountInMinorUnits)
	return NewMoneyFromMinorUnits(resultAmount, money.currency)
}

func (money Money) Add(
	other Money,
	conversionService CurrencyConversionService,
) (Money, error) {
	return money.operationWithNormalizedCurrency(
		other,
		conversionService,
		func(a int64, b int64) int64 {
			return a + b
		},
	)
}

func (money Money) Subtract(
	other Money,
	conversionService CurrencyConversionService,
) (Money, error) {
	return money.operationWithNormalizedCurrency(
		other,
		conversionService,
		func(a int64, b int64) int64 {
			return a - b
		},
	)
}

func (money Money) Multiply(factor int64) (Money, error) {
	if factor < 0 {
		return Money{}, fmt.Errorf("cannot multiply by negative factor")
	}

	resultAmount := money.amountInMinorUnits * factor
	return NewMoneyFromMinorUnits(resultAmount, money.currency)
}

func (money Money) normalizeCurrency(
	other Money,
	conversionService CurrencyConversionService,
) (Money, error) {
	if money.currency.Equals(other.currency) {
		return other, nil
	}

	convertedMoney, err := conversionService.Convert(other, money.currency)
	if err != nil {
		return Money{}, err
	}

	return convertedMoney, nil
}
