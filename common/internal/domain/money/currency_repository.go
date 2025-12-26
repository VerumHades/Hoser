package money

import "common/internal/shared"

type CurrencyRepository interface {
	Save(currency Currency) (Currency, error)

	FindByCode(code shared.CurrencyCode) (Currency, error)

	FetchNextBatch(
		lastSeenCurrencyCode shared.CurrencyCode,
		maximumBatchSize int,
	) ([]Currency, error)
}
