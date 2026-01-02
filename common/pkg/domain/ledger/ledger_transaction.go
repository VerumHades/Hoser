package ledger

import (
	"fmt"
	"time"

	"common/pkg/shared"
)

// ReferenceType defines the type of business event a ledger transaction represents.
type ReferenceType string

const (
	ReferenceTypePurchase   ReferenceType = "purchase"
	ReferenceTypeRefund     ReferenceType = "refund"
	ReferenceTypeChargeback ReferenceType = "chargeback"
	ReferenceTypePayout     ReferenceType = "payout"
	ReferenceTypeAdjustment ReferenceType = "adjustment"
)

// LedgerTransaction is the aggregate root representing a group of ledger entries
// for a single business event. It enforces the double-entry invariant.
type LedgerTransaction struct {
	id            shared.LedgerTransactionID
	referenceType ReferenceType
	referenceID   string
	createdAt     time.Time
	entries       []*LedgerEntry
}

func NewLedgerTransaction(
	referenceType ReferenceType,
	referenceID string,
	entries []*LedgerEntry,
) (*LedgerTransaction, error) {
	return NewLedgerTransactionWithID(
		shared.LedgerTransactionID(shared.GenerateUUID()),
		referenceType,
		referenceID,
		entries,
		time.Now(),
	)
}

// NewLedgerTransaction constructs a LedgerTransaction with its entries,
// enforcing the double-entry invariant (entries must sum to zero) and all other invariants.
func NewLedgerTransactionWithID(
	id shared.LedgerTransactionID,
	referenceType ReferenceType,
	referenceID string,
	entries []*LedgerEntry,
	createdAt time.Time,
) (*LedgerTransaction, error) {
	if id == "" {
		return nil, fmt.Errorf("ledger transaction ID cannot be empty")
	}
	if !isValidReferenceType(referenceType) {
		return nil, fmt.Errorf("invalid reference type: %s", referenceType)
	}
	if referenceID == "" {
		return nil, fmt.Errorf("reference ID cannot be empty")
	}
	if len(entries) < 2 {
		return nil, fmt.Errorf("ledger transaction must have at least two entries")
	}

	var total int64
	for _, e := range entries {
		total += e.AmountInMinorUnits()
	}
	if total != 0 {
		return nil, fmt.Errorf("ledger transaction entries must sum to zero, sum=%d", total)
	}

	// Defensive copy to preserve immutability
	copiedEntries := make([]*LedgerEntry, len(entries))
	copy(copiedEntries, entries)

	return &LedgerTransaction{
		id:            id,
		referenceType: referenceType,
		referenceID:   referenceID,
		createdAt:     createdAt,
		entries:       copiedEntries,
	}, nil
}

// Entries returns a read-only slice of ledger entries.
func (lt *LedgerTransaction) Entries() []*LedgerEntry {
	return lt.entries
}

// Getters for read-only access
func (lt *LedgerTransaction) ID() shared.LedgerTransactionID {
	return lt.id
}

func (lt *LedgerTransaction) ReferenceType() ReferenceType {
	return lt.referenceType
}

func (lt *LedgerTransaction) ReferenceID() string {
	return lt.referenceID
}

func (lt *LedgerTransaction) CreatedAt() time.Time {
	return lt.createdAt
}

// isValidReferenceType ensures the reference type is one of the allowed constants.
func isValidReferenceType(rt ReferenceType) bool {
	switch rt {
	case ReferenceTypePurchase, ReferenceTypeRefund, ReferenceTypeChargeback, ReferenceTypePayout, ReferenceTypeAdjustment:
		return true
	default:
		return false
	}
}
