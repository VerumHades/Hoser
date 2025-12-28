package user

import (
	"common/internal/domain/billing"
	"common/internal/domain/user"
	"common/internal/shared"
	"context"
	"fmt"
)

// PaymentProcessor defines the minimal interface to execute payments.
type PaymentProcessor interface {
	ProcessPayment(payment *billing.Payment, provider billing.PaymentProvider) (*billing.Payment, error)
}

type UserPaymentService struct {
	paymentRepository billing.PaymentCommandRepository
	payoutRepository  billing.PayoutCommandRepository

	accountRepository billing.BillingAccountQueryRepository
	userRepository    user.UserQueryRepository
	paymentProcessor  PaymentProcessor
}

func (service *UserPaymentService) Pay(
	ctx context.Context,
	billingAccountID shared.BillingAccountID,
	payment *billing.Payment,
) error {
	account, err := service.accountRepository.GetByID(ctx, billingAccountID)
	if err != nil {
		return fmt.Errorf("billing account not found: %w", err)
	}
	if !account.IsActive() {
		return fmt.Errorf("billing account is not active")
	}

	savedPayment, err := service.paymentRepository.Create(ctx, nil, payment)
	if err != nil {
		return err
	}

	_, err = service.paymentProcessor.ProcessPayment(savedPayment, account.PaymentProvider())
	if err != nil {
		return err
	}

	return nil
}

func (service *UserPaymentService) Payout(
	ctx context.Context,
	userID shared.UserID,
	payout *billing.Payout,
) error {
	_, err := service.userRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	_, err = service.payoutRepository.Create(ctx, nil, payout)
	if err != nil {
		return err
	}

	return nil
}
