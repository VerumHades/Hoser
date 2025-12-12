package payments

import (
	"common/pkg/billing/currency"
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
	Amount    currency.Money
	PaidAt    time.Time
}

type RefundReceipt struct {
	RefundID   string
	PaymentID  string
	Amount     currency.Money
	RefundedAt time.Time
}

type PaymentAuthorization struct {
	ID        string
	ExpiresAt time.Time
}

type PaymentGateway interface {
	AuthorizePayment(
		billingAccountID string,
		amount currency.Money,
	) (*PaymentAuthorization, error)

	CapturePayment(
		authorizationID string,
	) (*PaymentReceipt, error)

	RefundPayment(
		paymentID string,
		amount currency.Money,
	) (*RefundReceipt, error)
}
