package billing

import (
	"common/internal/shared"
	"time"
)

// PayoutRepository defines persistence operations for payouts.
type PayoutRepository interface {
	// Save persists a payout (insert or update).
	Save(payout *Payout) (*Payout, error)

	// GetByID retrieves a payout by its ID.
	GetByID(payoutID shared.PayoutID) (*Payout, error)

	// Delete removes a payout from storage.
	Delete(payoutID shared.PayoutID) error

	// FetchNextBatchByRecipient returns a batch of payouts for a recipient, ordered by creation date.
	// lastSeenPayoutID allows pagination to continue from the last record in previous batch.
	FetchNextBatchByRecipient(
		recipientID shared.UserID,
		lastSeenPayoutID shared.PayoutID,
		maximumBatchSize int,
	) ([]*Payout, error)

	// FetchNextBatchByStatus returns payouts filtered by status, paginated.
	FetchNextBatchByStatus(
		status PayoutStatus,
		lastSeenPayoutID shared.PayoutID,
		maximumBatchSize int,
	) ([]*Payout, error)

	// FetchNextBatchAvailableForClaim returns payouts that are claimable (status Available) in batches.
	FetchNextBatchAvailableForClaim(
		lastSeenPayoutID shared.PayoutID,
		maximumBatchSize int,
	) ([]*Payout, error)

	// FetchNextBatchCreatedBefore returns payouts created before a certain timestamp.
	FetchNextBatchCreatedBefore(
		cutoffTime time.Time,
		lastSeenPayoutID shared.PayoutID,
		maximumBatchSize int,
	) ([]*Payout, error)
}
