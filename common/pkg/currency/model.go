package currency

type Currency struct {
	Code   string // "USD"
	Name   string // "US Dollar"
	Symbol string // optional "$"
}

type Money struct {
	Amount   float64
	Currency Currency
}
