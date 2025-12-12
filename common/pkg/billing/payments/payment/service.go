package payment

import (
	"common/pkg/billing/currency"
	"time"
)

type PaymentService struct {
	repo PaymentRepository
}

func NewPaymentService(repo PaymentRepository) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) CreatePayment(billingAccountID string, amount currency.Money, metadata map[string]string) (*Payment, error) {
	payment := NewPendingPayment(billingAccountID, amount, metadata)
	err := s.repo.Save(payment)
	if err != nil {
		return nil, err
	}
	return payment, nil
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
