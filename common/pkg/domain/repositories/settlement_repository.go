package repositories

import (
	"common/pkg/domain/entities/accounting"
	"common/pkg/shared"
	"context"
)

// SettlementCommandRepository defines write operations for settlements.
type SettlementCommandRepository interface {
	// Create persists a new settlement within a transaction.
	Create(
		ctx context.Context,

		settlement *accounting.Settlement,
	) (*accounting.Settlement, error)

	// Update modifies an existing settlement within a transaction.
	Update(
		ctx context.Context,

		settlement *accounting.Settlement,
	) (*accounting.Settlement, error)

	// Delete removes a settlement by its ID within a transaction.
	Delete(
		ctx context.Context,

		settlementID shared.SettlementID,
	) error
}

// SettlementQueryRepository defines read-only operations for settlements.
type SettlementQueryRepository interface {
	// GetByID retrieves a settlement by its ID.
	GetByID(
		ctx context.Context,
		settlementID shared.SettlementID,
	) (*accounting.Settlement, error)

	// GetByLedgerTransactionID retrieves all settlements associated with a ledger transaction.
	GetByLedgerTransactionID(
		ctx context.Context,
		ledgerTransactionID shared.LedgerTransactionID,
	) ([]*accounting.Settlement, error)

	GetLastByLedgerTransactionID(
		ctx context.Context,
		ledgerTransactionID shared.LedgerTransactionID,
	) (*accounting.Settlement, error)
}
