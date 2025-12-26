package instance

import (
	"common/internal/domain/billing"
	"common/internal/domain/instance"
	"common/internal/domain/money"
	"common/internal/shared"

	"fmt"
	"time"
)

type PaymentService interface {
	HasUserBoughtListing(userID shared.UserID, listingID shared.ListingID) (bool, error)
}

// PaymentProcessor defines the minimal interface to execute payments.
type PaymentProcessor interface {
	ProcessPayment(payment *billing.Payment, provider billing.PaymentProvider) (*billing.Payment, error)
}

type HardwareCostCalculationService interface {
	CalculateCost(
		specification *shared.HardwareSpecification,
		duration time.Duration,
		at time.Time,
		currency shared.CurrencyCode,
	) (money.Money, error)
}

// InstanceSubscriptionService orchestrates instance creation, payment, and hardware cost calculations.
type InstanceSubscriptionService struct {
	paymentService                 PaymentService
	paymentProcessor               PaymentProcessor
	hardwareCostCalculationService HardwareCostCalculationService

	instanceRepository       instance.InstanceRepository
	paymentRepository        billing.PaymentRepository
	payoutRepository         billing.PayoutRepository
	billingAccountRepository billing.BillingAccountRepository

	currencyConcersionService money.CurrencyConversionService
}

func (s *InstanceSubscriptionService) CreateInstance(ownerID shared.UserID, listingID shared.ListingID, billingAccountID shared.BillingAccountID, hardwareSpec *shared.HardwareSpecification) (*instance.Instance, error) {
	hasBought, err := s.paymentService.HasUserBoughtListing(ownerID, listingID)
	if err != nil {
		return nil, err
	}
	if !hasBought {
		return nil, fmt.Errorf("cannot make instance of listing that hasnt been bought by the user")
	}

	instanceObj, err := instance.NewInstance(listingID, billingAccountID, hardwareSpec, time.Now().Add(-time.Hour))

	if err != nil {
		return nil, err
	}

	if err = instanceObj.Start(); err != nil {
		return nil, err
	}

	return s.instanceRepository.Save(instanceObj)
}

// RenewInstanceHardware renews an instance with a new hardware spec and/or duration.
// Charges only the difference between the already paid remaining period and the new total.
func (s *InstanceSubscriptionService) RenewInstanceHardware(
	instanceID shared.InstanceID,
	newHardwareSpec *shared.HardwareSpecification,
) (*instance.Instance, error) {

	currentInstance, err := s.instanceRepository.GetByID(instanceID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var remainingDuration time.Duration = currentInstance.RemainingDurationUntilExpiry()

	alreadyPaid, err := s.hardwareCostCalculationService.CalculateCost(
		currentInstance.HardwareSpecification(),
		remainingDuration,
		now,
		"USD",
	)
	if err != nil {
		return nil, err
	}

	newTotalCost, err := s.hardwareCostCalculationService.CalculateCost(
		newHardwareSpec,
		currentInstance.RenewalDuration(),
		now,
		"USD",
	)
	if err != nil {
		return nil, err
	}

	account, err := s.billingAccountRepository.GetByID(currentInstance.BillingAccountID())
	if err != nil {
		return nil, err
	}

	if alreadyPaid.AmountInMinorUnits() > newTotalCost.AmountInMinorUnits() {
		payoutAmount, err := alreadyPaid.Subtract(newTotalCost, s.currencyConcersionService)
		if err != nil {
			return nil, err
		}

		payout, err := billing.NewPayout(account.OwnerID(), payoutAmount)
		if err != nil {
			return nil, err
		}
		_, err = s.payoutRepository.Save(payout)
		if err != nil {
			return nil, err
		}
	} else if alreadyPaid.AmountInMinorUnits() < newTotalCost.AmountInMinorUnits() {
		amountToCharge, err := newTotalCost.Subtract(alreadyPaid, s.currencyConcersionService)

		if err != nil {
			return nil, err
		}

		payment, err := billing.NewSubscriptionPayment(
			currentInstance.BillingAccountID(),
			amountToCharge,
			instanceID,
			now.Add(currentInstance.RenewalDuration()),
			billing.SubscriptionPaymentInstanceHosting,
		)
		if err != nil {
			return nil, err
		}

		// Persist the payment before processing
		savedPayment, err := s.paymentRepository.Save(payment)
		if err != nil {
			return nil, err
		}

		// Execute the payment through the processor
		_, err = s.paymentProcessor.ProcessPayment(savedPayment, account.PaymentProvider())
		if err != nil {
			return nil, err
		}
	}

	currentInstance.Renew()
	currentInstance.SetHardwareSpecification(newHardwareSpec)

	updatedInstance, err := s.instanceRepository.Save(currentInstance)
	if err != nil {
		return nil, err
	}

	return updatedInstance, nil
}
