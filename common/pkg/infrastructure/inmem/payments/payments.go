package inmempayments

import (
	"common/pkg/billing/payments"
	"common/pkg/money"
	"errors"
	"sync"
	"time"
)

type InMemoryPaymentGateway struct {
	mu             sync.Mutex
	authorizations map[string]*payments.PaymentAuthorization
	payments       map[string]*payments.PaymentReceipt
	refunds        map[string]*payments.RefundReceipt
}

func NewInMemoryPaymentGateway() *InMemoryPaymentGateway {
	return &InMemoryPaymentGateway{
		authorizations: make(map[string]*payments.PaymentAuthorization),
		payments:       make(map[string]*payments.PaymentReceipt),
		refunds:        make(map[string]*payments.RefundReceipt),
	}
}

// AuthorizePayment creates a dummy payment authorization that expires in 5 minutes.
func (g *InMemoryPaymentGateway) AuthorizePayment(billingAccountID string, amount money.Money) (*payments.PaymentAuthorization, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	authID := "auth_" + time.Now().Format("20060102150405.000000000")
	auth := &payments.PaymentAuthorization{
		ID:        authID,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	g.authorizations[authID] = auth
	return auth, nil
}

// CapturePayment captures a previously authorized payment.
func (g *InMemoryPaymentGateway) CapturePayment(authorizationID string) (*payments.PaymentReceipt, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	auth, exists := g.authorizations[authorizationID]
	if !exists {
		return nil, errors.New("authorization not found")
	}
	if time.Now().After(auth.ExpiresAt) {
		return nil, errors.New("authorization expired")
	}

	paymentID := "pay_" + time.Now().Format("20060102150405.000000000")
	receipt := &payments.PaymentReceipt{
		PaymentID: paymentID,
		Amount:    money.Money{Amount: 1000, CurrencyCode: "USD"}, // Dummy amount
		PaidAt:    time.Now(),
	}
	g.payments[paymentID] = receipt

	// Remove authorization after capture
	delete(g.authorizations, authorizationID)

	return receipt, nil
}

// RefundPayment creates a dummy refund for a captured payment.
func (g *InMemoryPaymentGateway) RefundPayment(paymentID string, amount money.Money) (*payments.RefundReceipt, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	_, exists := g.payments[paymentID]
	if !exists {
		return nil, errors.New("payment not found")
	}

	refundID := "refund_" + time.Now().Format("20060102150405.000000000")
	refund := &payments.RefundReceipt{
		RefundID:   refundID,
		PaymentID:  paymentID,
		Amount:     amount,
		RefundedAt: time.Now(),
	}
	g.refunds[refundID] = refund
	return refund, nil
}
