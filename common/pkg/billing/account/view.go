package account

import (
	"common/pkg/billing/payments"
	"time"
)

type BillingAccountView struct {
	ID                string
	OwnerID           string
	Status            BillingAccountStatus
	PaymentProvider   payments.PaymentProvider
	ProviderAccountID string
	CreatedAt         time.Time
}

func (b *BillingAccount) ToView() *BillingAccountView {
	return &BillingAccountView{
		ID:                b.id,
		OwnerID:           b.ownerID,
		Status:            b.status,
		PaymentProvider:   b.paymentProvider,
		ProviderAccountID: b.providerAccountID,
		CreatedAt:         b.createdAt,
	}
}
