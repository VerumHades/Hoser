package ledger

import (
	"common/internal/shared"
	"context"
)

// LedgerTransactionCommandRepository defines write operations for ledger transactions.
type LedgerTransactionCommandRepository interface {
	// Create persists a new ledger transaction and its entries within a transaction.
	Create(
		ctx context.Context,
		transaction shared.Transaction,
		ledgerTransaction *LedgerTransaction,
	) error

	// Update allows modifications if absolutely necessary, e.g., metadata corrections.
	// In most systems, ledger transactions are immutable, so this might not be used.
	Update(
		ctx context.Context,
		transaction shared.Transaction,
		ledgerTransaction *LedgerTransaction,
	) error

	// Delete removes a ledger transaction within a transaction.
	// Usually discouraged in double-entry systems; only for special cases like testing or corrections.
	Delete(
		ctx context.Context,
		transaction shared.Transaction,
		ledgerTransactionID shared.LedgerTransactionID,
	) error
}

// LedgerTransactionQueryRepository defines read-only operations for ledger transactions.
type LedgerTransactionQueryRepository interface {
	// GetByID retrieves a ledger transaction by its ID.
	GetByID(
		ctx context.Context,
		ledgerTransactionID shared.LedgerTransactionID,
	) (*LedgerTransaction, error)

	// GetByReference retrieves all ledger transactions for a given business entity.
	GetByReference(
		ctx context.Context,
		referenceType ReferenceType,
		referenceID string,
	) ([]*LedgerTransaction, error)

	// GetByReference retrieves all ledger transactions for a given business entity.
	GetLatestByReferenceAndAccount(
		ctx context.Context,
		accountID shared.AccountID,
		referenceID string,
	) (*LedgerTransaction, error)

	// ListByAccount retrieves all ledger transactions affecting a specific account.
	ListByAccount(
		ctx context.Context,
		accountID shared.AccountID,
	) ([]*LedgerTransaction, error)
}
