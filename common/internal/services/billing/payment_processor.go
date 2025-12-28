package application

import (
	"common/internal/domain/billing"
	"context"
)

// PaymentProcessor executes one-time or subscription payments.
type PaymentProcessor struct {
	paymentRepository billing.PaymentCommandRepository
	gatewayRouter     billing.PaymentGatewayResolver
}

// NewPaymentProcessor constructs a PaymentProcessor.
func NewPaymentProcessor(
	paymentRepository billing.PaymentCommandRepository,
	gatewayRouter billing.PaymentGatewayResolver,
) *PaymentProcessor {
	return &PaymentProcessor{
		paymentRepository: paymentRepository,
		gatewayRouter:     gatewayRouter,
	}
}

// ProcessPayment executes the payment through the provider, marks status, and persists it.
func (p *PaymentProcessor) ProcessPayment(ctx context.Context, payment *billing.Payment, provider billing.PaymentProvider) (*billing.Payment, error) {
	gateway, err := p.gatewayRouter.Resolve(provider)
	if err != nil {
		payment.MarkFailed()
		_, _ = p.paymentRepository.Update(ctx, nil, payment)
		return payment, err
	}

	auth, err := gateway.AuthorizePayment(string(provider), payment.Amount())
	if err != nil {
		payment.MarkFailed()
		_, _ = p.paymentRepository.Update(ctx, nil, payment)
		return payment, err
	}

	if _, err := gateway.CapturePayment(auth.ID); err != nil {
		payment.MarkFailed()
		_, _ = p.paymentRepository.Update(ctx, nil, payment)
		return payment, err
	}

	payment.MarkPaid()
	savedPayment, err := p.paymentRepository.Update(ctx, nil, payment)
	if err != nil {
		return savedPayment, err
	}

	return savedPayment, nil
}
