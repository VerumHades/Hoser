package domain

import (
	"common/internal/shared"
	"time"
)

type PaymentGatewayResolver interface {
	Resolve(provider PaymentProvider) (PaymentGateway, error)
}

type PaymentProvider string

const (
	PaymentProviderStripe  PaymentProvider = "stripe"
	PaymentProviderInvoice PaymentProvider = "invoice"
)

type PaymentReceipt struct {
	PaymentID string
	Amount    shared.Money
	PaidAt    time.Time
}

type RefundReceipt struct {
	RefundID   string
	PaymentID  string
	Amount     shared.Money
	RefundedAt time.Time
}

type PaymentAuthorization struct {
	ID        string
	ExpiresAt time.Time
}

type PaymentGateway interface {
	AuthorizePayment(
		billingAccountID string,
		amount shared.Money,
	) (*PaymentAuthorization, error)

	CapturePayment(
		authorizationID string,
	) (*PaymentReceipt, error)

	RefundPayment(
		paymentID string,
		amount shared.Money,
	) (*RefundReceipt, error)
}
