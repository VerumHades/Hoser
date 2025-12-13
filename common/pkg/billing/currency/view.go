package currency

// CurrencyView is a read-only representation of a currency.
type CurrencyView struct {
	Code   string
	Name   string
	Symbol string
}

// ToView returns a read-only view of the Currency.
func (c Currency) ToView() CurrencyView {
	return CurrencyView{
		Code:   c.Code,
		Name:   c.Name,
		Symbol: c.Symbol,
	}
}
