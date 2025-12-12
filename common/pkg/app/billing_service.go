package app

import (
	"fmt"

	"common/pkg/billing/account"
	"common/pkg/billing/currency"
	"common/pkg/billing/payments"
	"common/pkg/billing/payments/payment"
)

type PaymentService struct {
	paymentService        *payment.PaymentService
	billingAccountService *account.BillingAccountService
	gatewayRouter         payments.PaymentGatewayResolver
}

// NewPaymentService creates a new payment service instance
func NewPaymentService(
	paymentService *payment.PaymentService,
	billingAccountService *account.BillingAccountService,
	gatewayRouter payments.PaymentGatewayResolver,
) *PaymentService {
	return &PaymentService{
		paymentService:        paymentService,
		billingAccountService: billingAccountService,
		gatewayRouter:         gatewayRouter,
	}
}

// CreateBillingAccount creates a new billing account for a user and returns a read-only view.
func (s *PaymentService) CreateBillingAccount(
	ownerID string,
	paymentProvider payments.PaymentProvider,
	providerAccountID string,
) (*account.BillingAccountView, error) {
	acct, err := s.billingAccountService.CreateAccount(ownerID, paymentProvider, providerAccountID)
	if err != nil {
		return nil, err
	}
	return acct.ToView(), nil
}

// Pay initiates a payment for the given billing account and amount.
func (s *PaymentService) Pay(
	billingAccountID string,
	amount currency.Money,
	metadata map[string]string,
) (*payment.PaymentView, error) {
	// Lookup billing account via the account service
	acct, err := s.billingAccountService.GetAccountView(billingAccountID)
	if err != nil {
		return nil, err
	}

	// Check account status
	if acct.Status == account.BillingAccountStatusSuspended || acct.Status == account.BillingAccountStatusClosed {
		return nil, fmt.Errorf("billing account is not active")
	}

	provider := acct.PaymentProvider

	// Create a pending payment using the payment service
	pmt, err := s.paymentService.CreatePayment(billingAccountID, amount, metadata)
	if err != nil {
		return nil, err
	}

	paymentView := pmt.ToView()

	// Resolve the gateway
	gateway, err := s.gatewayRouter.Resolve(provider)
	if err != nil {
		_ = s.paymentService.MarkPaymentFailed(paymentView.ID)
		return nil, err
	}

	// Authorize payment
	auth, err := gateway.AuthorizePayment(billingAccountID, amount)
	if err != nil {
		_ = s.paymentService.MarkPaymentFailed(paymentView.ID)
		return paymentView, err
	}

	// Capture payment
	if _, err = gateway.CapturePayment(auth.ID); err != nil {
		_ = s.paymentService.MarkPaymentFailed(paymentView.ID)
		return paymentView, err
	}

	// Mark payment as paid
	if err := s.paymentService.MarkPaymentPaid(paymentView.ID); err != nil {
		return paymentView, err
	}

	return paymentView, nil
}

// GetPaymentView returns a read-only view of the payment
func (s *PaymentService) GetPaymentView(paymentID string) (*payment.PaymentView, error) {
	return s.paymentService.GetPaymentView(paymentID)
}

// SuspendBillingAccount suspends a billing account
func (s *PaymentService) SuspendBillingAccount(accountID string) error {
	return s.billingAccountService.SuspendAccount(accountID)
}

// CloseBillingAccount closes a billing account
func (s *PaymentService) CloseBillingAccount(accountID string) error {
	return s.billingAccountService.CloseAccount(accountID)
}
