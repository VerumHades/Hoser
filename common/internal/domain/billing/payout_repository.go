package billing

import (
	"common/internal/shared"
	"context"
	"time"
)

// PayoutCommandRepository defines persistence operations that modify payouts.
type PayoutCommandRepository interface {
	// Create persists a new payout within a transaction.
	Create(
		ctx context.Context,
		transaction shared.Transaction,
		payout *Payout,
	) (*Payout, error)

	// Update modifies an existing payout within a transaction.
	Update(
		ctx context.Context,
		transaction shared.Transaction,
		payout *Payout,
	) (*Payout, error)

	// Delete removes a payout by its ID within a transaction.
	Delete(
		ctx context.Context,
		transaction shared.Transaction,
		payoutID shared.PayoutID,
	) error
}

// PayoutQueryRepository defines read-only operations for payouts.
type PayoutQueryRepository interface {
	// GetByID retrieves a payout by its ID.
	GetByID(
		ctx context.Context,
		payoutID shared.PayoutID,
	) (*Payout, error)

	// FetchNextBatchByRecipient returns a batch of payouts for a recipient.
	FetchNextBatchByRecipient(
		ctx context.Context,
		recipientID shared.UserID,
		request shared.BatchRequest,
	) (payouts []*Payout, nextCursor shared.Cursor, err error)

	// FetchNextBatchByStatus returns payouts filtered by status.
	FetchNextBatchByStatus(
		ctx context.Context,
		status PayoutStatus,
		request shared.BatchRequest,
	) (payouts []*Payout, nextCursor shared.Cursor, err error)

	// FetchNextBatchAvailableForClaim returns claimable payouts (status Available).
	FetchNextBatchAvailableForClaim(
		ctx context.Context,
		request shared.BatchRequest,
	) (payouts []*Payout, nextCursor shared.Cursor, er error)

	// FetchNextBatchCreatedBefore returns payouts created before a cutoff timestamp.
	FetchNextBatchCreatedBefore(
		ctx context.Context,
		cutoffTime time.Time,
		request shared.BatchRequest,
	) (payouts []*Payout, nextCursor shared.Cursor, er error)
}
