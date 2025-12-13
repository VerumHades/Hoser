package payment

import (
	"common/pkg/money"
	"time"
)

type PaymentView struct {
	ID               string
	BillingAccountID string
	Amount           money.Money
	Status           PaymentStatus
	CreatedAt        time.Time
	Kind             PaymentKind
	PaidAt           *time.Time
}

func (p *Payment) ToView() *PaymentView {
	return &PaymentView{
		ID:               p.id,
		BillingAccountID: p.billingAccountID,
		Amount:           p.amount,
		Status:           p.status,
		CreatedAt:        p.createdAt,
		PaidAt:           p.paidAt,
		Kind:             p.kind,
	}
}
