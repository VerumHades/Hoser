package app

import (
	"common/pkg/billing/account"
	"common/pkg/billing/payments/payment"
	"common/pkg/hardware"
	"common/pkg/instance"
)

// ReconciliationHandler defines actions triggered after instance changes.
type ReconciliationHandler interface {
	OnInstanceCreated(instance *instance.Instance) error
	OnInstanceUpdated(instance *instance.Instance) error
}

// InstanceEngineService orchestrates instance creation, payment, and reconciliation.
type InstanceEngineService struct {
	instanceService       *instance.InstanceService
	billingService        *account.BillingAccountService
	paymentService        *PaymentService
	publicListingService  *PublicListingService
	reconciliationHandler ReconciliationHandler
}

// NewInstanceEngineService creates a new instance engine service.
func NewInstanceEngineService(
	instanceService *instance.InstanceService,
	billingService *account.BillingAccountService,
	paymentService *PaymentService,
	publicListingService *PublicListingService,
	reconciliationHandler ReconciliationHandler,
) *InstanceEngineService {
	return &InstanceEngineService{
		instanceService:       instanceService,
		billingService:        billingService,
		paymentService:        paymentService,
		publicListingService:  publicListingService,
		reconciliationHandler: reconciliationHandler,
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

	// Trigger reconciliation if a handler is set
	if s.reconciliationHandler != nil {
		if err := s.reconciliationHandler.OnInstanceCreated(instanceObj); err != nil {
			return instanceObj, err
		}
	}

	return instanceObj, nil
}

// UpdateHardwareSpecification updates the hardware spec of an instance and triggers reconciliation.
func (s *InstanceEngineService) UpdateHardwareSpecification(instanceID string, newHardwareSpec *hardware.HardwareSpecification) (*instance.Instance, error) {
	instance, err := s.instanceService.UpdateHardwareSpecification(instanceID, newHardwareSpec)
	if err != nil {
		return nil, err
	}

	if s.reconciliationHandler != nil {
		if err := s.reconciliationHandler.OnInstanceUpdated(instance); err != nil {
			return instance, err
		}
	}

	return instance, nil
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
