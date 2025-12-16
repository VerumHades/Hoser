package app

import (
	"time"

	"common/pkg/billing/payments/payment"
	"common/pkg/hardware"
	"common/pkg/listing"
	"common/pkg/money"
)

type InstanceBillingService struct {
	listingService          *listing.ListingService
	paymentService          *PaymentService
	hardwareCostCalculation *HardwareCostCalculationService
}

// NewListingPurchaseService creates a new purchase service instance
func NewListingPurchaseService(
	listingService *listing.ListingService,
	paymentService *PaymentService,
	hardwareCostCalculation *HardwareCostCalculationService,
	conversionService money.CurrencyConversionService,
) *InstanceBillingService {
	return &InstanceBillingService{
		listingService:          listingService,
		paymentService:          paymentService,
		hardwareCostCalculation: hardwareCostCalculation,
	}
}

// BillInstance purchases a listing, taking into account past payments, hardware cost, and duration.
func (s *InstanceBillingService) BillInstance(
	billingAccountID string,
	listingID string,
	instanceID string,
	pricingID string,
	buyerUserID string,
	rentalDuration time.Duration,
	hardwareSpecification *hardware.HardwareSpecification,
) error {

	listingView, err := s.listingService.GetListing(listingID)
	if err != nil {
		return err
	}

	price := listingView.Price

	alreadyBought, err := s.paymentService.HasUserBoughtListing(buyerUserID, listingID)
	if err != nil {
		return err
	}

	if !alreadyBought {
		if _, err := s.paymentService.PayOneTime(billingAccountID, price, payment.OneTimePaymentMetadata{
			ListingID: listingID,
			UserID:    buyerUserID,
		}); err != nil {
			return err
		}
	}

	hardwareCost, err := s.hardwareCostCalculation.CalculateCost(
		hardwareSpecification,
		rentalDuration,
		time.Now(),
		price.CurrencyCode,
	)

	if err != nil {
		return err
	}

	if _, err := s.paymentService.PaySubscription(billingAccountID, hardwareCost, payment.SubscriptionPaymentMetadata{
		InstanceID:       instanceID,
		CurrentPeriodEnd: time.Now().Add(rentalDuration),
		Type:             payment.SubscriptionPaymentInstanceHosting,
	}); err != nil {
		return err
	}

	return nil
}
