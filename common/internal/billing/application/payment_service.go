package application

import (
	"fmt"
	"time"

	"common/internal/billing/domain"
	"common/internal/shared"
)

// PaymentService coordinates payments through domain factories and repositories.
type PaymentService struct {
	billingAccountRepository domain.BillingAccountRepository
	paymentRepository        domain.PaymentRepository
	gatewayRouter            domain.PaymentGatewayResolver
}

// NewPaymentService creates a new PaymentService instance.
func NewPaymentService(
	billingAccountRepository domain.BillingAccountRepository,
	paymentRepository domain.PaymentRepository,
	gatewayRouter domain.PaymentGatewayResolver,
) *PaymentService {
	return &PaymentService{
		billingAccountRepository: billingAccountRepository,
		paymentRepository:        paymentRepository,
		gatewayRouter:            gatewayRouter,
	}
}

// CreateBillingAccount creates a new billing account.
func (s *PaymentService) CreateBillingAccount(
	ownerID shared.UserID,
	paymentProvider domain.PaymentProvider,
	providerAccountID string,
) (*domain.BillingAccount, error) {
	account, err := domain.NewBillingAccount(ownerID, paymentProvider, providerAccountID)
	if err != nil {
		return nil, err
	}
	if err = s.billingAccountRepository.Save(account); err != nil {
		return nil, err
	}
	return account, nil
}

// PayOneTime executes a one-time payment.
func (s *PaymentService) PayOneTime(
	billingAccountID shared.BillingAccountID,
	amount shared.Money,
	listingID shared.ListingID,
	userID shared.UserID,
) (*domain.Payment, error) {
	account, err := s.billingAccountRepository.GetByID(billingAccountID)
	if err != nil {
		return nil, err
	}
	if !account.IsActive() {
		return nil, fmt.Errorf("billing account is not active")
	}

	payment, err := domain.NewOneTimePayment(billingAccountID, amount, listingID, userID)
	if err != nil {
		return nil, err
	}

	if err := s.paymentRepository.Save(payment); err != nil {
		return nil, err
	}

	if err := s.processPayment(payment, account.PaymentProvider()); err != nil {
		return payment, err
	}

	return payment, nil
}

// PaySubscription executes a subscription payment.
func (s *PaymentService) PaySubscription(
	billingAccountID shared.BillingAccountID,
	amount shared.Money,
	instanceID shared.InstanceID,
	periodEnd time.Time,
	paymentType domain.SubscriptionPaymentType,
) (*domain.Payment, error) {
	account, err := s.billingAccountRepository.GetByID(billingAccountID)
	if err != nil {
		return nil, err
	}
	if !account.IsActive() {
		return nil, fmt.Errorf("billing account is not active")
	}

	payment, err := domain.NewSubscriptionPayment(billingAccountID, amount, instanceID, periodEnd, paymentType)
	if err != nil {
		return nil, err
	}

	if err := s.paymentRepository.Save(payment); err != nil {
		return nil, err
	}

	if err := s.processPayment(payment, account.PaymentProvider()); err != nil {
		return payment, err
	}

	return payment, nil
}

// HasUserBoughtListing checks if a user has purchased a listing.
func (s *PaymentService) HasUserBoughtListing(userID shared.UserID, listingID shared.ListingID) (bool, error) {
	accounts, err := s.billingAccountRepository.ListByOwner(userID)
	if err != nil {
		return false, err
	}

	for _, account := range accounts {
		payments, err := s.paymentRepository.ListByBillingAccount(account.ID())
		if err != nil {
			return false, err
		}
		for _, payment := range payments {
			if payment.OneTimeMetadata() != nil &&
				payment.OneTimeMetadata().UserID() == userID &&
				payment.OneTimeMetadata().ListingID() == listingID {
				return true, nil
			}
		}
	}
	return false, nil
}

// SuspendBillingAccount suspends an account.
func (s *PaymentService) SuspendBillingAccount(accountID shared.BillingAccountID) error {
	account, err := s.billingAccountRepository.GetByID(accountID)
	if err != nil {
		return err
	}
	account.Suspend()
	return s.billingAccountRepository.Save(account)
}

// CloseBillingAccount closes an account.
func (s *PaymentService) CloseBillingAccount(accountID shared.BillingAccountID) error {
	account, err := s.billingAccountRepository.GetByID(accountID)
	if err != nil {
		return err
	}
	account.Close()
	return s.billingAccountRepository.Save(account)
}

// processPayment handles authorization, capture, and marking the payment status.
func (s *PaymentService) processPayment(payment *domain.Payment, provider domain.PaymentProvider) error {
	gateway, err := s.gatewayRouter.Resolve(provider)
	if err != nil {
		payment.MarkFailed()
		_ = s.paymentRepository.Save(payment)
		return err
	}

	auth, err := gateway.AuthorizePayment(string(provider), payment.Amount())
	if err != nil {
		payment.MarkFailed()
		_ = s.paymentRepository.Save(payment)
		return err
	}

	if _, err := gateway.CapturePayment(auth.ID); err != nil {
		payment.MarkFailed()
		_ = s.paymentRepository.Save(payment)
		return err
	}

	payment.MarkPaid()
	return s.paymentRepository.Save(payment)
}
