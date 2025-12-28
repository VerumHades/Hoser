package money

import (
	"common/internal/shared"
	"context"
)

// CurrencyCommandRepository defines write operations for currencies.
type CurrencyCommandRepository interface {
	// Create persists a new currency within a transaction.
	Create(
		ctx context.Context,
		transaction shared.Transaction,
		currency Currency,
	) (Currency, error)

	// Update modifies an existing currency within a transaction.
	Update(
		ctx context.Context,
		transaction shared.Transaction,
		currency Currency,
	) (Currency, error)
}

// CurrencyQueryRepository defines read-only operations for currencies.
type CurrencyQueryRepository interface {
	// FindByCode retrieves a currency by its code.
	FindByCode(
		ctx context.Context,
		code shared.CurrencyCode,
	) (Currency, error)

	// FetchNextBatch returns currencies in batches.
	FetchNextBatch(
		ctx context.Context,
		request shared.BatchRequest,
	) (currencies []Currency, nextCursor shared.Cursor, err error)
}
