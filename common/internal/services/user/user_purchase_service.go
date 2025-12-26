package user

import (
	"fmt"

	"common/internal/domain/billing"
	"common/internal/domain/money"
	"common/internal/shared"
)

// PaymentProcessor defines the minimal interface to execute payments.
type PaymentProcessor interface {
	ProcessPayment(payment *billing.Payment, provider billing.PaymentProvider) (*billing.Payment, error)
}

// UserPurchaseService handles all user listing purchase logic.
type UserPurchaseService struct {
	billingAccountRepository billing.BillingAccountRepository
	paymentRepository        billing.PaymentRepository
	paymentProcessor         PaymentProcessor
}

// NewUserPurchaseService constructs a new UserPurchaseService.
func NewUserPurchaseService(
	billingAccountRepository billing.BillingAccountRepository,
	paymentRepository billing.PaymentRepository,
	paymentProcessor PaymentProcessor,
) *UserPurchaseService {
	return &UserPurchaseService{
		billingAccountRepository: billingAccountRepository,
		paymentRepository:        paymentRepository,
		paymentProcessor:         paymentProcessor,
	}
}

// HasUserBoughtListing checks if a user already owns a listing.
func (s *UserPurchaseService) HasUserBoughtListing(userID shared.UserID, listingID shared.ListingID) (bool, error) {
	lastSeenAccountID := shared.BillingAccountID("")
	batchSize := 100

	for {
		accounts, err := s.billingAccountRepository.FetchNextBatchByOwner(userID, lastSeenAccountID, batchSize)
		if err != nil {
			return false, err
		}
		if len(accounts) == 0 {
			break
		}

		for _, account := range accounts {
			lastSeenAccountID = account.ID()
			lastSeenPaymentID := shared.PaymentID("")

			for {
				payments, err := s.paymentRepository.FetchNextBatchByBillingAccount(account.ID(), lastSeenPaymentID, batchSize)
				if err != nil {
					return false, err
				}
				if len(payments) == 0 {
					break
				}

				for _, payment := range payments {
					lastSeenPaymentID = payment.ID()
					if payment.OneTimeMetadata() != nil &&
						payment.OneTimeMetadata().UserID() == userID &&
						payment.OneTimeMetadata().ListingID() == listingID {
						return true, nil
					}
				}
			}
		}
	}

	return false, nil
}

// PurchaseListing ensures a user owns a listing and executes the payment if necessary.
func (s *UserPurchaseService) PurchaseListing(
	userID shared.UserID,
	listingID shared.ListingID,
	billingAccountID shared.BillingAccountID,
	amount money.Money,
) (*billing.Payment, error) {

	// Check ownership internally
	hasBought, err := s.HasUserBoughtListing(userID, listingID)
	if err != nil {
		return nil, fmt.Errorf("failed to check ownership: %w", err)
	}
	if hasBought {
		return nil, fmt.Errorf("user %s has already purchased listing %s", userID, listingID)
	}

	// Fetch the billing account
	account, err := s.billingAccountRepository.GetByID(billingAccountID)
	if err != nil {
		return nil, fmt.Errorf("billing account not found: %w", err)
	}
	if !account.IsActive() {
		return nil, fmt.Errorf("billing account is not active")
	}

	// Create the payment entity
	payment, err := billing.NewOneTimePayment(billingAccountID, amount, listingID, userID)
	if err != nil {
		return nil, err
	}

	// Persist the payment before processing
	savedPayment, err := s.paymentRepository.Save(payment)
	if err != nil {
		return nil, err
	}

	// Execute the payment through the processor
	finalPayment, err := s.paymentProcessor.ProcessPayment(savedPayment, account.PaymentProvider())
	if err != nil {
		return finalPayment, err
	}

	return finalPayment, nil
}
