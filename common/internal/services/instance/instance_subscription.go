package instance

import (
	"common/internal/domain/billing"
	"common/internal/domain/instance"
	"common/internal/domain/money"
	"common/internal/shared"
	"context"

	"fmt"
	"time"
)

type OwnershipQueryService interface {
	HasUserBoughtListing(ctx context.Context, userID shared.UserID, listingID shared.ListingID) (bool, error)
}

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
	paymentService                 OwnershipQueryService
	paymentProvider                PaymentProvider
	hardwareCostCalculationService HardwareCostCalculationService

	instanceCommandRepository instance.InstanceCommandRepository
	instanceQueryRepository   instance.InstanceQueryRepository

	billingAccountRepository billing.BillingAccountQueryRepository

	currencyConcersionService money.CurrencyConversionService
}

func (s *InstanceSubscriptionService) CreateInstance(ctx context.Context, ownerID shared.UserID, listingID shared.ListingID, billingAccountID shared.BillingAccountID, hardwareSpec *shared.HardwareSpecification) (*instance.Instance, error) {
	hasBought, err := s.paymentService.HasUserBoughtListing(ctx, ownerID, listingID)
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

	return s.instanceCommandRepository.Create(ctx, nil, instanceObj)
}

// RenewInstanceHardware renews an instance with a new hardware spec and/or duration.
// Charges only the difference between the already paid remaining period and the new total.
func (s *InstanceSubscriptionService) RenewInstanceHardware(
	ctx context.Context,
	instanceID shared.InstanceID,
	newHardwareSpec *shared.HardwareSpecification,
) (*instance.Instance, error) {
	currentInstance, err := s.instanceQueryRepository.GetByID(ctx, instanceID)
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

		if err = s.paymentProvider.Payout(ctx, account.OwnerID(), payout); err != nil {
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

		if err = s.paymentProvider.Pay(ctx, currentInstance.BillingAccountID(), payment); err != nil {
			return nil, err
		}
	}

	currentInstance.Renew()
	currentInstance.SetHardwareSpecification(newHardwareSpec)

	updatedInstance, err := s.instanceCommandRepository.Update(ctx, nil, currentInstance)
	if err != nil {
		return nil, err
	}

	return updatedInstance, nil
}
