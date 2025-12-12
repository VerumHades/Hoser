package payment

import (
	"common/pkg/billing/currency"
	"time"
)

type PaymentView struct {
	ID               string
	BillingAccountID string
	Amount           currency.MoneyView
	Status           PaymentStatus
	CreatedAt        time.Time
	PaidAt           *time.Time
	Metadata         map[string]string
}

func (p *Payment) ToView() *PaymentView {
	return &PaymentView{
		ID:               p.id,
		BillingAccountID: p.billingAccountID,
		Amount:           p.amount.ToView(),
		Status:           p.status,
		CreatedAt:        p.createdAt,
		PaidAt:           p.paidAt,
		Metadata:         p.metadata,
	}
}
