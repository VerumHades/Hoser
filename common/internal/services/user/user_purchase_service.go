package user

import (
	"context"
	"fmt"

	"common/internal/domain/billing"
	"common/internal/domain/money"
	"common/internal/shared"
	"common/internal/util"
)

type PaymentProvider interface {
	Pay(
		ctx context.Context,
		billingAccountID shared.BillingAccountID,
		payment *billing.Payment,
	) error
	Payout(
		ctx context.Context,
		userID shared.UserID,
		payout *billing.Payout,
	) error
}

// UserPurchaseService handles all user listing purchase logic.
type UserPurchaseService struct {
	billingAccountRepository billing.BillingAccountQueryRepository
	paymentRepository        billing.PaymentQueryRepository

	paymentProvider PaymentProvider
}

// NewUserPurchaseService constructs a new UserPurchaseService.
func NewUserPurchaseService(
	billingAccountRepository billing.BillingAccountQueryRepository,
	paymentRepository billing.PaymentQueryRepository,
	paymentProvider PaymentProvider,
) *UserPurchaseService {
	return &UserPurchaseService{
		billingAccountRepository: billingAccountRepository,
		paymentRepository:        paymentRepository,
		paymentProvider:          paymentProvider,
	}
}

// HasUserBoughtListing checks if a user already owns a listing.
func (s *UserPurchaseService) HasUserBoughtListing(ctx context.Context, userID shared.UserID, listingID shared.ListingID) (bool, error) {
	// Look for a matching payment across all billing accounts
	foundPayment, err := util.LookupInBatches(
		ctx,
		100,
		func(ctx context.Context, request shared.BatchRequest) ([]*billing.BillingAccount, shared.Cursor, error) {
			return s.billingAccountRepository.FetchNextBatchByOwner(ctx, userID, request)
		},
		func(ctx context.Context, account *billing.BillingAccount) (*billing.Payment, error) {
			return util.LookupInBatches(
				ctx,
				100,
				func(ctx context.Context, request shared.BatchRequest) ([]*billing.Payment, shared.Cursor, error) {
					return s.paymentRepository.FetchNextBatchByBillingAccount(ctx, account.ID(), request)
				},
				func(ctx context.Context, payment *billing.Payment) (*billing.Payment, error) {
					metadata := payment.OneTimeMetadata()
					if metadata != nil &&
						metadata.UserID() == userID &&
						metadata.ListingID() == listingID {
						return payment, nil
					}
					return nil, nil
				},
			)
		},
	)

	if err != nil {
		return false, err
	}

	return foundPayment != nil, nil
}

// PurchaseListing ensures a user owns a listing and executes the payment if necessary.
func (s *UserPurchaseService) PurchaseListing(
	ctx context.Context,
	userID shared.UserID,
	listingID shared.ListingID,
	billingAccountID shared.BillingAccountID,
	amount money.Money,
) error {
	// Check ownership internally
	hasBought, err := s.HasUserBoughtListing(ctx, userID, listingID)
	if err != nil {
		return fmt.Errorf("failed to check ownership: %w", err)
	}
	if hasBought {
		return fmt.Errorf("user %s has already purchased listing %s", userID, listingID)
	}

	// Create the payment entity
	payment, err := billing.NewOneTimePayment(billingAccountID, amount, listingID, userID)
	if err != nil {
		return err
	}

	return s.paymentProvider.Pay(ctx, billingAccountID, payment)
}
