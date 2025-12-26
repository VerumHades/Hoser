package billing

import (
	"common/internal/shared"
	"time"
)

type PaymentRepository interface {
	Save(payment *Payment) (*Payment, error)

	GetByID(paymentID shared.PaymentID) (*Payment, error)

	FetchNextBatchByBillingAccount(
		accountID shared.BillingAccountID,
		lastSeenPaymentID shared.PaymentID,
		maximumBatchSize int,
	) ([]*Payment, error)

	FetchNextBatchByBillingAccountAndTimeRange(
		accountID shared.BillingAccountID,
		startTime time.Time,
		endTime time.Time,
		lastSeenPaymentID shared.PaymentID,
		maximumBatchSize int,
	) ([]*Payment, error)
}
