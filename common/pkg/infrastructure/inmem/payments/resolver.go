package inmempayments

import (
	"common/pkg/billing/payments"
	"errors"
)

// InMemoryPaymentGatewayResolver always resolves to InMemoryPaymentGateway
type InMemoryPaymentGatewayResolver struct {
	gateway *InMemoryPaymentGateway
}

// NewInMemoryPaymentGatewayResolver creates a new resolver instance
func NewInMemoryPaymentGatewayResolver() *InMemoryPaymentGatewayResolver {
	return &InMemoryPaymentGatewayResolver{
		gateway: NewInMemoryPaymentGateway(),
	}
}

// Resolve returns the in-memory gateway regardless of provider
func (r *InMemoryPaymentGatewayResolver) Resolve(provider payments.PaymentProvider) (payments.PaymentGateway, error) {
	if r.gateway == nil {
		return nil, errors.New("in-memory gateway not initialized")
	}
	return r.gateway, nil
}
