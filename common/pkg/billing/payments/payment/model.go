package payment

import (
	"time"

	"common/pkg/money"
	"common/pkg/util"
)

type PaymentKind string

const (
	PaymentTypeOneTime      PaymentKind = "one_time"
	PaymentTypeSubscription PaymentKind = "subscription"
	PaymentTypeUnknown      PaymentKind = "unknown"
)

type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusPaid     PaymentStatus = "paid"
	PaymentStatusFailed   PaymentStatus = "failed"
	PaymentStatusRefunded PaymentStatus = "refunded"
)

type Payment struct {
	id               string
	billingAccountID string
	amount           money.Money
	status           PaymentStatus
	kind             PaymentKind
	createdAt        time.Time
	paidAt           *time.Time
}

type OneTimePaymentMetadata struct {
	ListingID string
	UserID    string
}

type SubscriptionPaymentMetadata struct {
	InstanceID       string
	CurrentPeriodEnd time.Time
}

// NewPendingPayment creates a new payment with optional metadata
func NewPendingPayment(
	billingAccountID string,
	amount money.Money,
) *Payment {

	return &Payment{
		id:               util.GenerateUUID(),
		billingAccountID: billingAccountID,
		amount:           amount,
		status:           PaymentStatusPending,
		kind:             PaymentTypeUnknown,
		createdAt:        time.Now(),
	}
}
func (payment *Payment) MarkPaid(paidAt time.Time) {
	payment.status = PaymentStatusPaid
	payment.paidAt = &paidAt
}

func (payment *Payment) MarkFailed() {
	payment.status = PaymentStatusFailed
}

func (payment *Payment) MarkRefunded() {
	payment.status = PaymentStatusRefunded
}
