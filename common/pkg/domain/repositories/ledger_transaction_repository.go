package repositories

import (
	"common/pkg/domain/entities/accounting"
	"common/pkg/shared"
	"context"
)

// LedgerTransactionCommandRepository defines write operations for ledger transactions.
type LedgerTransactionCommandRepository interface {
	// Create persists a new ledger transaction and its entries within a transaction.
	Create(
		ctx context.Context,

		ledgerTransaction *accounting.LedgerTransaction,
	) error

	// Update allows modifications if absolutely necessary, e.g., metadata corrections.
	// In most systems, ledger transactions are immutable, so this might not be used.
	Update(
		ctx context.Context,

		ledgerTransaction *accounting.LedgerTransaction,
	) error

	// Delete removes a ledger transaction within a transaction.
	// Usually discouraged in double-entry systems; only for special cases like testing or corrections.
	Delete(
		ctx context.Context,

		ledgerTransactionID shared.LedgerTransactionID,
	) error
}

// LedgerTransactionQueryRepository defines read-only operations for ledger transactions.
type LedgerTransactionQueryRepository interface {
	// GetByID retrieves a ledger transaction by its ID.
	GetByID(
		ctx context.Context,
		ledgerTransactionID shared.LedgerTransactionID,
	) (*accounting.LedgerTransaction, error)

	// GetByReference retrieves all ledger transactions for a given business entity.
	GetByReference(
		ctx context.Context,
		referenceType accounting.ReferenceType,
		referenceID string,
	) ([]*accounting.LedgerTransaction, error)

	// GetByReference retrieves all ledger transactions for a given business entity.
	GetLatestByReferenceAndAccount(
		ctx context.Context,
		accountID shared.AccountID,
		referenceID string,
	) (*accounting.LedgerTransaction, error)
}
