package account

import (
	"time"

	"common/pkg/billing/payments"
	"common/pkg/util"
)

type BillingAccountStatus string

const (
	BillingAccountStatusActive    BillingAccountStatus = "active"
	BillingAccountStatusSuspended BillingAccountStatus = "suspended"
	BillingAccountStatusClosed    BillingAccountStatus = "closed"
)

type BillingAccount struct {
	id                string
	ownerID           string
	status            BillingAccountStatus
	paymentProvider   payments.PaymentProvider
	providerAccountID string
	createdAt         time.Time
}

func NewBillingAccount(
	ownerID string,
	paymentProvider payments.PaymentProvider,
	providerAccountID string,
) *BillingAccount {
	return &BillingAccount{
		id:                util.GenerateUUID(),
		ownerID:           ownerID,
		status:            BillingAccountStatusActive,
		paymentProvider:   paymentProvider,
		providerAccountID: providerAccountID,
		createdAt:         time.Now(),
	}
}

func (b *BillingAccount) IsChargeable() bool {
	return b.status == BillingAccountStatusActive
}

func (b *BillingAccount) Suspend() {
	b.status = BillingAccountStatusSuspended
}

func (b *BillingAccount) Close() {
	b.status = BillingAccountStatusClosed
}
