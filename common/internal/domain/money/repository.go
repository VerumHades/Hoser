package money

import "common/internal/shared"

type CurrencyRepository interface {
	FindByCode(code shared.CurrencyCode) (Currency, error)
	Save(currency Currency) error
	ListAll() ([]Currency, error)
}
