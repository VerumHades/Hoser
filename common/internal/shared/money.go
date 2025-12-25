package shared

import (
	"errors"
	"fmt"
)

type CurrencyConversionService interface {
	Convert(m Money, targetCurrency string) (Money, error)
}

type Money struct {
	amount       float64
	currencyCode string
}

// NewMoney creates a validated Money instance.
func NewMoney(amount float64, currencyCode string) (*Money, error) {
	if currencyCode == "" {
		return nil, errors.New("currency code cannot be empty")
	}
	if amount < 0 {
		return nil, errors.New("amount cannot be negative")
	}
	return &Money{amount: amount, currencyCode: currencyCode}, nil
}

// Accessors
func (m *Money) Amount() float64 {
	return m.amount
}

func (m *Money) CurrencyCode() string {
	return m.currencyCode
}

func (m *Money) IsZero() bool {
	return m.amount == 0
}

func (m *Money) Add(other *Money, conversionService CurrencyConversionService) (*Money, error) {
	otherBase := other
	if other.currencyCode != m.currencyCode {
		var err error
		converted, err := conversionService.Convert(*other, m.currencyCode)
		if err != nil {
			return nil, err
		}
		otherBase = &converted
	}

	resultAmount := m.amount + otherBase.amount
	return NewMoney(resultAmount, m.currencyCode)
}

func (m *Money) Subtract(other *Money, conversionService CurrencyConversionService) (*Money, error) {
	otherBase := other
	if other.currencyCode != m.currencyCode {
		var err error
		converted, err := conversionService.Convert(*other, m.currencyCode)
		if err != nil {
			return nil, err
		}
		otherBase = &converted
	}

	resultAmount := m.amount - otherBase.amount
	if resultAmount < 0 {
		return nil, fmt.Errorf("resulting amount cannot be negative")
	}

	return NewMoney(resultAmount, m.currencyCode)
}

func (m *Money) Multiply(factor float64) (*Money, error) {
	if factor < 0 {
		return nil, fmt.Errorf("cannot multiply by negative factor")
	}
	return NewMoney(m.amount*factor, m.currencyCode)
}
