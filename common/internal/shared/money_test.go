package shared_test

import (
	"errors"
	"math"
	"testing"

	"common/internal/shared"

	"github.com/stretchr/testify/assert"
)

type mockConversionService struct{}

func (m *mockConversionService) Convert(money shared.Money, targetCurrency string) (shared.Money, error) {
	if targetCurrency == "ERR" || money.CurrencyCode() == "ERR" {
		return shared.Money{}, errors.New("conversion error")
	}

	// Simulate conversion: double the amount and create via NewMoney to enforce domain rules
	converted, err := shared.NewMoney(money.Amount()*2, targetCurrency)
	if err != nil {
		return shared.Money{}, err
	}
	return *converted, nil
}
func TestNewMoney_Validation(t *testing.T) {
	// Valid money
	m, err := shared.NewMoney(10, "USD")
	assert.NoError(t, err)
	assert.Equal(t, 10.0, m.Amount())
	assert.Equal(t, "USD", m.CurrencyCode())

	// Negative amount
	_, err = shared.NewMoney(-5, "USD")
	assert.Error(t, err)

	// Empty currency
	_, err = shared.NewMoney(5, "")
	assert.Error(t, err)
}

func TestMoney_Add_EdgeCases(t *testing.T) {
	conversion := &mockConversionService{}

	m1, _ := shared.NewMoney(10, "USD")
	m2, _ := shared.NewMoney(5, "USD")
	mZero, _ := shared.NewMoney(0, "USD")
	mEUR, _ := shared.NewMoney(10, "EUR")
	mERR, _ := shared.NewMoney(5, "ERR")

	// Add zero
	result, err := m1.Add(mZero, conversion)
	assert.NoError(t, err)
	assert.Equal(t, 10.0, result.Amount())

	// Add same currency
	result, err = m1.Add(m2, conversion)
	assert.NoError(t, err)
	assert.Equal(t, 15.0, result.Amount())

	// Add different currency
	result, err = m1.Add(mEUR, conversion)
	assert.NoError(t, err)
	assert.Equal(t, 10+10*2, result.Amount()) // mEUR converted to USD

	// Conversion error
	_, err = m1.Add(mERR, conversion)
	assert.Error(t, err)

	// Very large amounts
	large1, _ := shared.NewMoney(math.MaxFloat64/2, "USD")
	large2, _ := shared.NewMoney(math.MaxFloat64/3, "USD")
	result, err = large1.Add(large2, conversion)
	assert.NoError(t, err)
	assert.Equal(t, math.MaxFloat64/2+math.MaxFloat64/3, result.Amount())
}

func TestMoney_Subtract_EdgeCases(t *testing.T) {
	conversion := &mockConversionService{}

	m1, _ := shared.NewMoney(10, "USD")
	m2, _ := shared.NewMoney(5, "USD")
	mEUR, _ := shared.NewMoney(5, "EUR")
	mERR, _ := shared.NewMoney(5, "ERR")

	// Subtract same currency
	result, err := m1.Subtract(m2, conversion)
	assert.NoError(t, err)
	assert.Equal(t, 5.0, result.Amount())

	// Subtract to zero
	result, err = m2.Subtract(m2, conversion)
	assert.NoError(t, err)
	assert.Equal(t, 0.0, result.Amount())

	// Subtract different currency
	result, err = m1.Subtract(mEUR, conversion)
	assert.NoError(t, err)
	assert.Equal(t, 10-5*2, result.Amount()) // mEUR converted

	// Subtract too large → should error (cannot be negative)
	_, err = m2.Subtract(m1, conversion)
	assert.Error(t, err)

	// Conversion error
	_, err = m1.Subtract(mERR, conversion)
	assert.Error(t, err)
}

func TestMoney_Multiply_EdgeCases(t *testing.T) {
	m, _ := shared.NewMoney(10, "USD")

	// Multiply by zero
	result, err := m.Multiply(0)
	assert.NoError(t, err)
	assert.Equal(t, 0.0, result.Amount())

	// Multiply by positive factor
	result, err = m.Multiply(2.5)
	assert.NoError(t, err)
	assert.Equal(t, 25.0, result.Amount())

	// Multiply by negative → should error
	_, err = m.Multiply(-2)
	assert.Error(t, err)

	// Multiply by very large factor
	result, err = m.Multiply(1e10)
	assert.NoError(t, err)
	assert.Equal(t, 10*1e10, result.Amount())
}
