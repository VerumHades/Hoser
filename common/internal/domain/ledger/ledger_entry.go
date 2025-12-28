package ledger

import (
	"common/internal/shared"
	"fmt"
	"time"
)

// LedgerEntry is a domain entity representing a debit or credit applied to an account.
type LedgerEntry struct {
	accountID          shared.AccountID
	amountInMinorUnits int64
	createdAt          time.Time
}

// NewLedgerEntry constructs a new ledger entry while enforcing invariants.
func NewLedgerEntry(id shared.LedgerEntryID, accountID shared.AccountID, amountInMinorUnits int64) (*LedgerEntry, error) {
	if id == "" {
		return nil, fmt.Errorf("ledger entry ID cannot be empty")
	}
	if accountID == "" {
		return nil, fmt.Errorf("account ID cannot be empty")
	}
	if amountInMinorUnits == 0 {
		return nil, fmt.Errorf("ledger entry amount cannot be zero")
	}

	return &LedgerEntry{
		accountID:          accountID,
		amountInMinorUnits: amountInMinorUnits,
		createdAt:          time.Now(),
	}, nil
}

func (e *LedgerEntry) AccountID() shared.AccountID {
	return e.accountID
}

func (e *LedgerEntry) AmountInMinorUnits() int64 {
	return e.amountInMinorUnits
}

func (e *LedgerEntry) CreatedAt() time.Time {
	return e.createdAt
}
