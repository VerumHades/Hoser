package billing

import (
	"common/internal/shared"
	"context"
	"time"
)

type PaymentCommandRepository interface {
	// Create persists a new payment record.
	Create(
		ctx context.Context,
		transaction shared.Transaction,
		payment *Payment,
	) (*Payment, error)

	// Update modifies an existing payment record.
	Update(
		ctx context.Context,
		transaction shared.Transaction,
		payment *Payment,
	) (*Payment, error)
}

type PaymentQueryRepository interface {
	GetByID(
		ctx context.Context,
		paymentID shared.PaymentID,
	) (*Payment, error)

	FetchNextBatchByBillingAccount(
		ctx context.Context,
		accountID shared.BillingAccountID,
		request shared.BatchRequest,
	) (payments []*Payment, nextCursor shared.Cursor, err error)

	FetchNextBatchByBillingAccountAndTimeRange(
		ctx context.Context,
		accountID shared.BillingAccountID,
		startTime time.Time,
		endTime time.Time,
		lastSeenPaymentID shared.PaymentID,
		request shared.BatchRequest,
	) (payments []*Payment, nextCursor shared.Cursor, err error)
}
