package app

import (
	"fmt"
	"time"

	"common/pkg/billing/payments/payment"
	"common/pkg/hardware"
	"common/pkg/listing"
	"common/pkg/money"
)

type ListingPurchaseService struct {
	listingService          *listing.ListingService
	paymentService          *PaymentService
	hardwareCostCalculation *HardwareCostCalculationService
	conversionService       money.CurrencyConversionService
}

// NewListingPurchaseService creates a new purchase service instance
func NewListingPurchaseService(
	listingService *listing.ListingService,
	paymentService *PaymentService,
	hardwareCostCalculation *HardwareCostCalculationService,
	conversionService money.CurrencyConversionService,
) *ListingPurchaseService {
	return &ListingPurchaseService{
		listingService:          listingService,
		paymentService:          paymentService,
		hardwareCostCalculation: hardwareCostCalculation,
		conversionService:       conversionService,
	}
}

// BillInstance purchases a listing, taking into account past payments, hardware cost, and duration.
func (s *ListingPurchaseService) BillInstance(
	billingAccountID string,
	listingID string,
	pricingID string,
	buyerUserID string,
	rentalDuration time.Duration,
	hardwareSpecification *hardware.HardwareSpecification,
) (*payment.PaymentView, error) {

	listingView, err := s.listingService.GetListingView(listingID)
	if err != nil {
		return nil, err
	}

	var selectedPricing *listing.PricingView
	for _, pricing := range listingView.Pricing {
		if pricing.ID == pricingID {
			selectedPricing = pricing
			break
		}
	}

	if selectedPricing == nil {
		return nil, fmt.Errorf("pricing not found")
	}

	alreadyBought, err := s.paymentService.HasUserBoughtListing(buyerUserID, listingID)
	if err != nil {
		return nil, err
	}

	totalAmount := selectedPricing.Amount
	if alreadyBought {
		totalAmount = money.Money{Amount: 0, CurrencyCode: "usd"}
	}

	hardwareCost, err := s.hardwareCostCalculation.CalculateCost(
		hardwareSpecification,
		rentalDuration,
		time.Now(),
		totalAmount.CurrencyCode,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate hardware cost: %w", err)
	}

	totalAmount, err = totalAmount.Add(hardwareCost, s.conversionService)
	if err != nil {
		return nil, fmt.Errorf("failed to add hardware cost: %w", err)
	}

	switch selectedPricing.Type {
	case listing.OneTime:
		return s.paymentService.PayOneTime(
			billingAccountID,
			totalAmount,
			payment.OneTimePaymentMetadata{
				ListingID: listingID,
				UserID:    buyerUserID,
			},
		)
	case listing.Monthly, listing.Yearly:
		return s.paymentService.PaySubscription(
			billingAccountID,
			totalAmount,
			payment.SubscriptionPaymentMetadata{
				InstanceID:       listingID,
				CurrentPeriodEnd: calculatePeriodEnd(selectedPricing.Type, rentalDuration),
			},
		)
	default:
		return nil, fmt.Errorf("unsupported pricing type")
	}
}

func calculatePeriodEnd(pricingType listing.PricingType, duration time.Duration) time.Time {
	now := time.Now()
	switch pricingType {
	case listing.Monthly:
		return now.Add(duration) // duration could be months converted properly
	case listing.Yearly:
		return now.Add(duration) // duration could be years converted properly
	default:
		return now
	}
}
