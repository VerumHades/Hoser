package payment

import (
	"common/pkg/money"
	"fmt"
	"time"
)

type PaymentService struct {
	repo PaymentRepository
}

func NewPaymentService(repo PaymentRepository) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) CreatePayment(billingAccountID string, amount money.Money) (*Payment, error) {
	payment := NewPendingPayment(billingAccountID, amount)
	err := s.repo.Save(payment)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (s *PaymentService) AddOneTimePaymentMetadata(
	paymentID string,
	metadata OneTimePaymentMetadata,
) error {
	payment, err := s.repo.GetByID(paymentID)
	if err != nil {
		return err
	}

	if payment.kind != PaymentTypeUnknown {
		return fmt.Errorf("cannot add one-time metadata: payment kind is already %s", payment.kind)
	}

	payment.kind = PaymentTypeOneTime

	if err := s.repo.AttachOneTimeMetadata(paymentID, metadata); err != nil {
		return err
	}

	return s.repo.Save(payment)
}

func (s *PaymentService) AddSubscriptionPaymentMetadata(
	paymentID string,
	metadata SubscriptionPaymentMetadata,
) error {
	payment, err := s.repo.GetByID(paymentID)
	if err != nil {
		return err
	}

	if payment.kind != PaymentTypeUnknown {
		return fmt.Errorf("cannot add subscription metadata: payment kind is already %s", payment.kind)
	}

	payment.kind = PaymentTypeSubscription

	if err := s.repo.AttachSubscriptionMetadata(paymentID, metadata); err != nil {
		return err
	}

	return s.repo.Save(payment)
}

func (s *PaymentService) GetOneTimePaymentMetadata(paymentID string) (OneTimePaymentMetadata, error) {
	payment, err := s.repo.GetByID(paymentID)
	if err != nil {
		return OneTimePaymentMetadata{}, err
	}

	if payment.kind != PaymentTypeOneTime {
		return OneTimePaymentMetadata{}, fmt.Errorf("payment does not have one-time metadata")
	}

	metadata, err := s.repo.GetOneTimeMetadata(paymentID)
	if err != nil {
		return OneTimePaymentMetadata{}, err
	}

	return metadata, nil
}

func (s *PaymentService) GetSubscriptionPaymentMetadata(paymentID string) (SubscriptionPaymentMetadata, error) {
	payment, err := s.repo.GetByID(paymentID)
	if err != nil {
		return SubscriptionPaymentMetadata{}, err
	}

	if payment.kind != PaymentTypeSubscription {
		return SubscriptionPaymentMetadata{}, fmt.Errorf("payment does not have subscription metadata")
	}

	metadata, err := s.repo.GetSubscriptionMetadata(paymentID)
	if err != nil {
		return SubscriptionPaymentMetadata{}, err
	}

	return metadata, nil
}

func (s *PaymentService) ListByBillingAccount(accountID string) ([]*PaymentView, error) {
	payments, err := s.repo.ListByBillingAccount(accountID)
	if err != nil {
		return nil, err
	}

	views := make([]*PaymentView, len(payments))
	for i, p := range payments {
		views[i] = p.ToView()
	}

	return views, nil
}

func (s *PaymentService) ListAllPaymentsByBillingAccount(accountID string) ([]*PaymentView, error) {
	payments, err := s.repo.ListByBillingAccount(accountID)
	if err != nil {
		return nil, err
	}

	views := make([]*PaymentView, len(payments))
	for i, p := range payments {
		views[i] = p.ToView()
	}
	return views, nil
}

func (s *PaymentService) MarkPaymentPaid(paymentID string) error {
	payment, err := s.repo.GetByID(paymentID)
	if err != nil {
		return err
	}
	payment.MarkPaid(time.Now())
	return s.repo.Save(payment)
}

func (s *PaymentService) MarkPaymentFailed(paymentID string) error {
	payment, err := s.repo.GetByID(paymentID)
	if err != nil {
		return err
	}
	payment.MarkFailed()
	return s.repo.Save(payment)
}

func (s *PaymentService) GetPaymentView(paymentID string) (*PaymentView, error) {
	payment, err := s.repo.GetByID(paymentID)
	if err != nil {
		return nil, err
	}
	return payment.ToView(), nil
}
