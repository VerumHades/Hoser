package accounting

import (
	"time"

	"common/pkg/shared"
)

/**
 * LedgerTransactionCreatedEvent is emitted when a new ledger transaction
 * is successfully initialized and satisfies the double-entry invariant.
 */
type LedgerTransactionCreatedEvent struct {
	ID            shared.LedgerTransactionID
	ReferenceType ReferenceType
	ReferenceID   string
	CreatedAt     time.Time
}

/**
 * LedgerTransactionUpdateEvent is scheduled for future use when
 * transaction metadata might need modification.
 */
type LedgerTransactionUpdateEvent struct {
	ID            shared.LedgerTransactionID
	ReferenceType *ReferenceType
	ReferenceID   *string
}
