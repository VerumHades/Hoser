package accounting

import (
	"fmt"
	"time"

	"common/pkg/domain/entities/events"
	"common/pkg/shared"
)

/**
 * ReferenceType defines the type of business event a ledger transaction represents.
 */
type ReferenceType string

const (
	ReferenceTypePurchase   ReferenceType = "purchase"
	ReferenceTypeRefund     ReferenceType = "refund"
	ReferenceTypeChargeback ReferenceType = "chargeback"
	ReferenceTypePayout     ReferenceType = "payout"
	ReferenceTypeAdjustment ReferenceType = "adjustment"
)

/**
 * LedgerTransaction is the aggregate root representing a group of ledger entries
 * for a single business event. It enforces the double-entry invariant.
 */
type LedgerTransaction struct {
	id            shared.LedgerTransactionID
	referenceType ReferenceType
	referenceID   string
	createdAt     time.Time
	entries       []*LedgerEntry
}

/**
 * NewLedgerTransaction constructs a LedgerTransaction with its entries,
 * enforcing the double-entry invariant and returning a creation domain event.
 */
func NewLedgerTransaction(
	referenceType ReferenceType,
	referenceID string,
	ledgerEntries []*LedgerEntry,
) (*LedgerTransaction, events.DomainEventEnvelope[LedgerTransactionCreatedEvent], error) {
	return NewLedgerTransactionWithID(
		shared.LedgerTransactionID(shared.GenerateUUID()),
		referenceType,
		referenceID,
		ledgerEntries,
		time.Now(),
	)
}

/**
 * NewLedgerTransactionWithID constructs a LedgerTransaction with a specific ID and timestamp,
 * validating all invariants before emitting the domain event.
 */
func NewLedgerTransactionWithID(
	ledgerTransactionID shared.LedgerTransactionID,
	referenceType ReferenceType,
	referenceID string,
	ledgerEntries []*LedgerEntry,
	createdAt time.Time,
) (*LedgerTransaction, events.DomainEventEnvelope[LedgerTransactionCreatedEvent], error) {
	if ledgerTransactionID == "" {
		return nil, events.DomainEventEnvelope[LedgerTransactionCreatedEvent]{}, fmt.Errorf("ledger transaction ID cannot be empty")
	}
	if !isValidReferenceType(referenceType) {
		return nil, events.DomainEventEnvelope[LedgerTransactionCreatedEvent]{}, fmt.Errorf("invalid reference type: %s", referenceType)
	}
	if referenceID == "" {
		return nil, events.DomainEventEnvelope[LedgerTransactionCreatedEvent]{}, fmt.Errorf("reference ID cannot be empty")
	}
	if len(ledgerEntries) < 2 {
		return nil, events.DomainEventEnvelope[LedgerTransactionCreatedEvent]{}, fmt.Errorf("ledger transaction must have at least two entries")
	}

	var totalBalance int64
	for _, entry := range ledgerEntries {
		totalBalance += entry.AmountInMinorUnits()
	}

	if totalBalance != 0 {
		return nil, events.DomainEventEnvelope[LedgerTransactionCreatedEvent]{}, fmt.Errorf("ledger transaction entries must sum to zero, sum=%d", totalBalance)
	}

	copiedEntries := make([]*LedgerEntry, len(ledgerEntries))
	copy(copiedEntries, ledgerEntries)

	ledgerTransaction := &LedgerTransaction{
		id:            ledgerTransactionID,
		referenceType: referenceType,
		referenceID:   referenceID,
		createdAt:     createdAt,
		entries:       copiedEntries,
	}

	eventPayload := LedgerTransactionCreatedEvent{
		ID:            ledgerTransactionID,
		ReferenceType: referenceType,
		ReferenceID:   referenceID,
		CreatedAt:     createdAt,
	}

	return ledgerTransaction, events.NewDomainEventEnvelope(eventPayload), nil
}

/**
 * Entries returns a read-only slice of ledger entries.
 */
func (lt *LedgerTransaction) Entries() []*LedgerEntry {
	return lt.entries
}

/**
 * ID returns the unique identifier for the transaction.
 */
func (lt *LedgerTransaction) ID() shared.LedgerTransactionID {
	return lt.id
}

/**
 * ReferenceType returns the business category of the transaction.
 */
func (lt *LedgerTransaction) ReferenceType() ReferenceType {
	return lt.referenceType
}

/**
 * ReferenceID returns the external identifier associated with this transaction.
 */
func (lt *LedgerTransaction) ReferenceID() string {
	return lt.referenceID
}

/**
 * CreatedAt returns the timestamp when the transaction was initialized.
 */
func (lt *LedgerTransaction) CreatedAt() time.Time {
	return lt.createdAt
}

/**
 * isValidReferenceType ensures the reference type is one of the allowed constants.
 */
func isValidReferenceType(referenceType ReferenceType) bool {
	switch referenceType {
	case ReferenceTypePurchase, ReferenceTypeRefund, ReferenceTypeChargeback, ReferenceTypePayout, ReferenceTypeAdjustment:
		return true
	default:
		return false
	}
}
