package payment

import (
	"time"

	"common/pkg/billing/currency"
	"common/pkg/util"
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
	amount           currency.Money
	status           PaymentStatus
	createdAt        time.Time
	paidAt           *time.Time
	metadata         map[string]string
}

// NewPendingPayment creates a new payment with optional metadata
func NewPendingPayment(billingAccountID string, amount currency.Money, metadata map[string]string) *Payment {
	return &Payment{
		id:               util.GenerateUUID(),
		billingAccountID: billingAccountID,
		amount:           amount,
		status:           PaymentStatusPending,
		createdAt:        time.Now(),
		metadata:         metadata,
	}
}

func (p *Payment) MarkPaid(paidAt time.Time) {
	p.status = PaymentStatusPaid
	p.paidAt = &paidAt
}

func (p *Payment) MarkFailed() {
	p.status = PaymentStatusFailed
}

func (p *Payment) MarkRefunded() {
	p.status = PaymentStatusRefunded
}
