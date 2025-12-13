package money

type CurrencyConversionService interface {
	Convert(m Money, targetCurrency string) (Money, error)
}

type Money struct {
	Amount       float64
	CurrencyCode string
}

// Add converts 'other' to the base currency and then adds it to 'm'.
// Returns the sum in the base currency.
func (m Money) Add(other Money, conversionService CurrencyConversionService) (Money, error) {
	baseCurrency := m.CurrencyCode

	mBase := m
	if m.CurrencyCode != baseCurrency {
		var err error
		mBase, err = conversionService.Convert(m, baseCurrency)
		if err != nil {
			return Money{}, err
		}
	}

	otherBase := other
	if other.CurrencyCode != baseCurrency {
		var err error
		otherBase, err = conversionService.Convert(other, baseCurrency)
		if err != nil {
			return Money{}, err
		}
	}

	// Add amounts
	return Money{
		Amount:       mBase.Amount + otherBase.Amount,
		CurrencyCode: baseCurrency,
	}, nil
}

// Multiply multiplies the amount by a factor and returns a new Money value.
func (m Money) Multiply(factor float64) Money {
	return Money{
		Amount:       m.Amount * factor,
		CurrencyCode: m.CurrencyCode,
	}
}
