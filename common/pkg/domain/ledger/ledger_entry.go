package ledger

import (
	"common/pkg/shared"
	"fmt"
	"time"
)

// LedgerEntry is a domain entity representing a debit or credit applied to an account.
type LedgerEntry struct {
	accountID          shared.AccountID
	amountInMinorUnits int64
	createdAt          time.Time
}

type LedgerEntryPair struct {
	accountID shared.AccountID
	amount    int64
}

func LedgerEntryShorthand(accountID shared.AccountID, amount int64) LedgerEntryPair {
	return LedgerEntryPair{
		accountID: accountID,
		amount:    amount,
	}
}

func NewLedgerEntries(pairs ...struct {
	accountID shared.AccountID
	amount    int64
}) ([]*LedgerEntry, error) {
	entries := make([]*LedgerEntry, 0, len(pairs))
	for _, p := range pairs {
		e, err := NewLedgerEntry(p.accountID, p.amount)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// ReverseLedgerEntries returns a new slice of ledger entries that exactly
// cancels out the effect of the given entries. Useful for refunds or rollbacks.
func ReverseLedgerEntries(entries []*LedgerEntry) ([]*LedgerEntry, error) {
	reversed := make([]*LedgerEntry, 0, len(entries))

	for _, e := range entries {
		reversedEntry, err := NewLedgerEntry(e.AccountID(), -e.AmountInMinorUnits())
		if err != nil {
			return nil, fmt.Errorf("failed to create reversed ledger entry for %s: %w", e.AccountID(), err)
		}
		reversed = append(reversed, reversedEntry)
	}

	return reversed, nil
}

// NewLedgerEntry constructs a new ledger entry while enforcing invariants.
func NewLedgerEntry(accountID shared.AccountID, amountInMinorUnits int64) (*LedgerEntry, error) {
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
