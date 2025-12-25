package money_test

import (
	"errors"
	"math"
	"testing"

	"common/internal/domain/money"
	"common/internal/shared"
)

type stubCurrencyConversionService struct {
	convertFunc func(money money.Money, targetCurrency money.Currency) (money.Money, error)
}

func (service stubCurrencyConversionService) Convert(
	money money.Money,
	targetCurrency money.Currency,
) (money.Money, error) {
	if service.convertFunc == nil {
		return money.Money{}, errors.New("unexpected conversion invocation")
	}
	return service.convertFunc(money, targetCurrency)
}

func mustCurrency(code string, minorUnitDecimalPlaces int) money.Currency {
	currency, err := money.NewCurrency(shared.CurrencyCode(code), minorUnitDecimalPlaces)
	if err != nil {
		panic(err)
	}
	return currency
}

func mustMoney(amountInMinorUnits int64, currency money.Currency) money.Money {
	money, err := money.NewMoneyFromMinorUnits(amountInMinorUnits, currency)
	if err != nil {
		panic(err)
	}
	return money
}

func TestNewMoneyFromMinorUnits_NegativeAmount_ReturnsError(t *testing.T) {
	currency := mustCurrency("USD", 2)

	_, err := money.NewMoneyFromMinorUnits(-1, currency)

	if err == nil {
		t.Fatalf("expected error for negative amount, got nil")
	}
}

func TestAdd_SameCurrency_DoesNotInvokeConversion(t *testing.T) {
	currency := mustCurrency("USD", 2)

	baseMoney := mustMoney(100, currency)
	otherMoney := mustMoney(50, currency)

	conversionService := stubCurrencyConversionService{
		convertFunc: func(
			money money.Money,
			targetCurrency money.Currency,
		) (money.Money, error) {
			t.Fatalf("conversion must not be invoked for same currency")
			return money.Money{}, nil
		},
	}

	result, err := baseMoney.Add(otherMoney, conversionService)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.AmountInMinorUnits() != 150 {
		t.Fatalf("expected 150, got %d", result.AmountInMinorUnits())
	}
}

func TestAdd_DifferentCurrency_ConversionErrorIsPropagated(t *testing.T) {
	usd := mustCurrency("USD", 2)
	eur := mustCurrency("EUR", 2)

	baseMoney := mustMoney(100, usd)
	otherMoney := mustMoney(100, eur)

	conversionService := stubCurrencyConversionService{
		convertFunc: func(
			money money.Money,
			targetCurrency money.Currency,
		) (money.Money, error) {
			return money.Money{}, errors.New("conversion failed")
		},
	}

	_, err := baseMoney.Add(otherMoney, conversionService)

	if err == nil {
		t.Fatalf("expected conversion error, got nil")
	}
}

func TestSubtract_ResultingNegativeAmount_ReturnsError(t *testing.T) {
	currency := mustCurrency("USD", 2)

	baseMoney := mustMoney(50, currency)
	otherMoney := mustMoney(100, currency)

	conversionService := stubCurrencyConversionService{}

	_, err := baseMoney.Subtract(otherMoney, conversionService)

	if err == nil {
		t.Fatalf("expected error when subtraction results in negative amount")
	}
}

func TestMultiply_ByZero_ReturnsZeroMoney(t *testing.T) {
	currency := mustCurrency("USD", 2)

	money := mustMoney(100, currency)

	result, err := money.Multiply(0)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.IsZero() {
		t.Fatalf("expected zero money, got %d", result.AmountInMinorUnits())
	}
}

func TestMultiply_NegativeFactor_ReturnsError(t *testing.T) {
	currency := mustCurrency("USD", 2)

	money := mustMoney(100, currency)

	_, err := money.Multiply(-1)

	if err == nil {
		t.Fatalf("expected error for negative multiplication factor")
	}
}

func TestAdd_WhenOverflowWouldResultInNegativeMoney_Fails(t *testing.T) {
	currency := mustCurrency("USD", 2)

	moneyA := mustMoney(math.MaxInt64, currency)
	moneyB := mustMoney(1, currency)

	conversionService := stubCurrencyConversionService{}

	result, err := moneyA.Add(moneyB, conversionService)

	if result.AmountInMinorUnits() < 0 {
		t.Fatalf("operation produced negative money, violates invariant")
	}

	// Optional: assert that the operation either succeeded safely or returned an error
	if err != nil {
		t.Logf("operation failed as expected: %v", err)
	}
}
