package ledger

import (
	"fmt"
	"time"

	"common/internal/shared"
)

// SettlementStatus defines the possible states of a settlement.
type SettlementStatus string

const (
	SettlementStatusPending   SettlementStatus = "pending"
	SettlementStatusCompleted SettlementStatus = "completed"
	SettlementStatusFailed    SettlementStatus = "failed"
)

// Settlement represents the state of an external payment or settlement associated
// with a ledger transaction or business event.
type Settlement struct {
	id                  shared.SettlementID
	ledgerTransactionID shared.LedgerTransactionID
	accountID           shared.AccountID
	amountInMinorUnits  int64
	status              SettlementStatus
	referenceID         string
	createdAt           time.Time
	updatedAt           time.Time
}

// NewSettlement constructs a new settlement with initial status "pending".
func NewSettlement(
	ledgerTransactionID shared.LedgerTransactionID,
	accountID shared.AccountID,
	amountInMinorUnits int64,
	referenceID string,
) (*Settlement, error) {
	if ledgerTransactionID == "" {
		return nil, fmt.Errorf("ledger transaction ID cannot be empty")
	}
	if accountID == "" {
		return nil, fmt.Errorf("account ID cannot be empty")
	}
	if amountInMinorUnits <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}
	if referenceID == "" {
		return nil, fmt.Errorf("reference ID cannot be empty")
	}

	return &Settlement{
		id:                  shared.SettlementID(shared.GenerateUUID()),
		ledgerTransactionID: ledgerTransactionID,
		accountID:           accountID,
		amountInMinorUnits:  amountInMinorUnits,
		status:              SettlementStatusPending,
		referenceID:         referenceID,
		createdAt:           time.Now(),
		updatedAt:           time.Now(),
	}, nil
}

// MarkCompleted sets the settlement status to completed.
func (s *Settlement) MarkCompleted() {
	s.status = SettlementStatusCompleted
	s.updatedAt = time.Now()
}

// MarkFailed sets the settlement status to failed.
func (s *Settlement) MarkFailed() {
	s.status = SettlementStatusFailed
	s.updatedAt = time.Now()
}

// Getters for read-only access
func (s *Settlement) ID() shared.SettlementID {
	return s.id
}

func (s *Settlement) LedgerTransactionID() shared.LedgerTransactionID {
	return s.ledgerTransactionID
}

func (s *Settlement) AccountID() shared.AccountID {
	return s.accountID
}

func (s *Settlement) Amount() int64 {
	return s.amountInMinorUnits
}

func (s *Settlement) Status() SettlementStatus {
	return s.status
}

func (s *Settlement) ReferenceID() string {
	return s.referenceID
}

func (s *Settlement) CreatedAt() time.Time {
	return s.createdAt
}

func (s *Settlement) UpdatedAt() time.Time {
	return s.updatedAt
}
