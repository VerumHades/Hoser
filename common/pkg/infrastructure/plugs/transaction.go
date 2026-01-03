package plugs

import (
	"context"
)

// FakeTransactionProvider is a no-op transaction provider.
type FakeTransactionProvider struct{}

// NewFakeTransactionProvider creates a new FakeTransactionProvider.
func NewFakeTransactionProvider() *FakeTransactionProvider {
	return &FakeTransactionProvider{}
}

// BeginTransaction returns a no-op transaction that satisfies shared.Transaction.
func (provider *FakeTransactionProvider) BeginTransaction(ctx context.Context) (context.Context, error) {
	return ctx, nil
}
