package app

import (
	"common/pkg/billing/account"
	"common/pkg/billing/payments/payment"
	"common/pkg/hardware"
	"common/pkg/instance"
	"common/pkg/money"
	"time"
)

// InstanceEngineService orchestrates instance creation, payment, and hardware cost calculations.
type InstanceEngineService struct {
	instanceService                *instance.InstanceService
	billingService                 *account.BillingAccountService
	paymentService                 *PaymentService
	publicListingService           *PublicListingService
	hardwareCostCalculationService *HardwareCostCalculationService
}

// NewInstanceEngineService creates a new instance engine service with the calculation service.
func NewInstanceEngineService(
	instanceService *instance.InstanceService,
	billingService *account.BillingAccountService,
	paymentService *PaymentService,
	publicListingService *PublicListingService,
	hardwareCostCalculationService *HardwareCostCalculationService,
) *InstanceEngineService {
	return &InstanceEngineService{
		instanceService:                instanceService,
		billingService:                 billingService,
		paymentService:                 paymentService,
		publicListingService:           publicListingService,
		hardwareCostCalculationService: hardwareCostCalculationService,
	}
}

// LaunchInstance launches a new instance for a user and automatically purchases the listing if not already owned.
func (s *InstanceEngineService) LaunchInstance(ownerID, listingID, billingAccountID string, hardwareSpec *hardware.HardwareSpecification) (*instance.Instance, error) {
	// Fetch the listing to get its price
	listingView, err := s.publicListingService.GetPublicListing(listingID)
	if err != nil {
		return nil, err
	}

	// Check if the user already bought the listing
	hasBought, err := s.paymentService.HasUserBoughtListing(ownerID, listingID)
	if err != nil {
		return nil, err
	}

	// If not purchased, perform a one-time payment using the listing price
	if !hasBought {
		paymentMetadata := payment.OneTimePaymentMetadata{
			UserID:    ownerID,
			ListingID: listingID,
		}

		if _, err := s.paymentService.PayOneTime(billingAccountID, listingView.Price, paymentMetadata); err != nil {
			return nil, err
		}
	}

	// Create the instance
	instanceObj, err := s.instanceService.CreateInstance(listingID, billingAccountID, hardwareSpec)
	if err != nil {
		return nil, err
	}

	return instanceObj, nil
}

// RenewInstanceHardware renews an instance with a new hardware spec and/or duration.
// Charges only the difference between the already paid remaining period and the new total.
func (s *InstanceEngineService) RenewInstanceHardware(
	instanceID string,
	newHardwareSpec *hardware.HardwareSpecification,
	duration time.Duration,
) (*instance.Instance, error) {

	// Fetch the instance
	currentInstance, err := s.instanceService.GetInstance(instanceID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var remainingDuration time.Duration
	if currentInstance.Expiry.After(now) {
		remainingDuration = currentInstance.Expiry.Sub(now)
	} else {
		remainingDuration = 0
	}

	// Cost already paid for remaining time with current spec
	alreadyPaid, err := s.hardwareCostCalculationService.CalculateCost(
		currentInstance.HardwareSpecification,
		remainingDuration,
		now,
		"USD",
	)
	if err != nil {
		return nil, err
	}

	// Total cost for new spec over total duration
	newTotalCost, err := s.hardwareCostCalculationService.CalculateCost(
		newHardwareSpec,
		duration,
		now,
		"USD",
	)
	if err != nil {
		return nil, err
	}

	amountToCharge, err := newTotalCost.Subtract(alreadyPaid, s.hardwareCostCalculationService.conversionService)
	if err != nil {
		return nil, err
	}
	if amountToCharge.Amount < 0 {
		amountToCharge = money.Money{Amount: 0, CurrencyCode: "USD"} // no refunds on downgrade
	}

	expiry := now.Add(duration)

	if amountToCharge.Amount > 0 {
		_, err := s.paymentService.PaySubscription(currentInstance.BillingID, amountToCharge, payment.SubscriptionPaymentMetadata{
			InstanceID:       instanceID,
			CurrentPeriodEnd: expiry,
			Type:             payment.SubscriptionPaymentInstanceHosting,
		})
		if err != nil {
			return nil, err
		}
	}

	// Update instance spec and expiry
	updatedInstance, err := s.instanceService.UpdateHardwareSpecification(instanceID, newHardwareSpec)
	if err != nil {
		return nil, err
	}
	updatedInstance, err = s.instanceService.SetExpiry(instanceID, expiry)
	if err != nil {
		return nil, err
	}

	return updatedInstance, nil
}

// ListInstancesByBillingAccount lists all instances for a billing account.
func (s *InstanceEngineService) ListInstancesByBillingAccount(billingID string) ([]*instance.Instance, error) {
	return s.instanceService.ListByBillingAccount(billingID)
}

// GetInstance retrieves an instance by ID.
func (s *InstanceEngineService) GetInstance(instanceID string) (*instance.Instance, error) {
	return s.instanceService.GetInstance(instanceID)
}

// ListInstancesByOwner lists all instances for a user by resolving billing accounts.
func (s *InstanceEngineService) ListInstancesByOwner(ownerID string) ([]*instance.Instance, error) {
	accounts, err := s.billingService.ListByOwner(ownerID)
	if err != nil {
		return nil, err
	}

	var instances []*instance.Instance
	for _, account := range accounts {
		accountInstances, err := s.instanceService.ListByBillingAccount(account.ID)
		if err != nil {
			return nil, err
		}
		instances = append(instances, accountInstances...)
	}
	return instances, nil
}
