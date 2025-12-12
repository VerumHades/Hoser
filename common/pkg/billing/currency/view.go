package currency

// CurrencyView is a read-only representation of a currency.
type CurrencyView struct {
	Code   string
	Name   string
	Symbol string
}

// MoneyView is a read-only representation of a money amount.
type MoneyView struct {
	Amount   float64
	Currency CurrencyView
}

// ToView returns a read-only view of the Currency.
func (c Currency) ToView() CurrencyView {
	return CurrencyView{
		Code:   c.Code,
		Name:   c.Name,
		Symbol: c.Symbol,
	}
}

// ToView returns a read-only view of the Money.
func (m Money) ToView() MoneyView {
	return MoneyView{
		Amount:   m.Amount,
		Currency: m.Currency.ToView(),
	}
}
