package plugs

import (
	"context"

	"common/pkg/shared"
)

// FakeTransaction is a no-op transaction used for standalone MongoDB or testing.
type FakeTransaction struct{}

// SessionContext just returns the original context; no session is used.
func (f *FakeTransaction) SessionContext(ctx context.Context) context.Context {
	return ctx
}

// Commit does nothing and always succeeds.
func (f *FakeTransaction) Commit() error {
	return nil
}

// Rollback does nothing and always succeeds.
func (f *FakeTransaction) Rollback() error {
	return nil
}

// FakeTransactionProvider is a no-op transaction provider.
type FakeTransactionProvider struct{}

// NewFakeTransactionProvider creates a new FakeTransactionProvider.
func NewFakeTransactionProvider() *FakeTransactionProvider {
	return &FakeTransactionProvider{}
}

// BeginTransaction returns a no-op transaction that satisfies shared.Transaction.
func (provider *FakeTransactionProvider) BeginTransaction(ctx context.Context) (shared.Transaction, error) {
	return &FakeTransaction{}, nil
}

// Ensure FakeTransaction satisfies shared.Transaction
var _ shared.Transaction = (*FakeTransaction)(nil)
