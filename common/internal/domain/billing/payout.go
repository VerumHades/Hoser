package billing

import (
	"errors"
	"time"

	"common/internal/domain/money"
	"common/internal/shared"
)

// PayoutStatus represents the current state of a payout.
type PayoutStatus int

const (
	PayoutPending    PayoutStatus = iota // not yet claimable
	PayoutAvailable                      // can be claimed
	PayoutProcessing                     // in the process of being paid
	PayoutPaid                           // completed
	PayoutFailed                         // failed to process
)

// Payout represents an amount that a user can claim or withdraw.
type Payout struct {
	id          shared.PayoutID
	recipientID shared.UserID // who can claim this payout
	amount      money.Money
	status      PayoutStatus

	createdAt   time.Time
	claimedAt   *time.Time // timestamp when payout was claimed
	completedAt *time.Time // timestamp when payout was actually paid
}

// NewPayout creates a new payout with initial status Pending.
func NewPayout(recipientID shared.UserID, amount money.Money) (*Payout, error) {
	if recipientID == "" {
		return nil, errors.New("recipient ID cannot be empty")
	}
	if amount.IsZero() {
		return nil, errors.New("payout amount cannot be zero")
	}

	return &Payout{
		id:          shared.PayoutID(shared.GenerateUUID()),
		recipientID: recipientID,
		amount:      amount,
		status:      PayoutPending,
		createdAt:   time.Now().UTC(),
	}, nil
}

// ID returns the payout ID.
func (p *Payout) ID() shared.PayoutID {
	return p.id
}

// RecipientID returns the ID of the user who can claim this payout.
func (p *Payout) RecipientID() shared.UserID {
	return p.recipientID
}

// Amount returns the payout amount.
func (p *Payout) Amount() money.Money {
	return p.amount
}

// Status returns the current payout status.
func (p *Payout) Status() PayoutStatus {
	return p.status
}

// CreatedAt returns the creation timestamp of the payout.
func (p *Payout) CreatedAt() time.Time {
	return p.createdAt
}

// ClaimedAt returns the timestamp when the payout was claimed, if any.
func (p *Payout) ClaimedAt() *time.Time {
	return p.claimedAt
}

// CompletedAt returns the timestamp when the payout was completed, if any.
func (p *Payout) CompletedAt() *time.Time {
	return p.completedAt
}

// MarkAvailable sets the payout as claimable.
func (p *Payout) MarkAvailable() error {
	if p.status != PayoutPending {
		return errors.New("payout can only become available from Pending state")
	}
	p.status = PayoutAvailable
	return nil
}

// Claim marks the payout as claimed.
func (p *Payout) Claim() error {
	if p.status != PayoutAvailable {
		return errors.New("payout can only be claimed from Available state")
	}
	now := time.Now().UTC()
	p.claimedAt = &now
	p.status = PayoutProcessing
	return nil
}

// Complete marks the payout as completed/paid.
func (p *Payout) Complete() error {
	if p.status != PayoutProcessing {
		return errors.New("payout can only be completed from Processing state")
	}
	now := time.Now().UTC()
	p.completedAt = &now
	p.status = PayoutPaid
	return nil
}

// Fail marks the payout as failed.
func (p *Payout) Fail() error {
	if p.status != PayoutProcessing && p.status != PayoutAvailable {
		return errors.New("payout can only fail from Available or Processing state")
	}
	p.status = PayoutFailed
	return nil
}

// Validate ensures the payout entity is in a consistent state.
func (p *Payout) Validate() error {
	if p.id == "" {
		return errors.New("payout ID cannot be empty")
	}
	if p.recipientID == "" {
		return errors.New("recipient ID cannot be empty")
	}
	if p.amount.IsZero() {
		return errors.New("payout amount cannot be zero")
	}
	switch p.status {
	case PayoutPending, PayoutAvailable, PayoutProcessing, PayoutPaid, PayoutFailed:
		// valid
	default:
		return errors.New("invalid payout status")
	}
	return nil
}
