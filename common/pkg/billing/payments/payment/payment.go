package payment

import (
	"fmt"
	"time"

	"common/pkg/money"
	"common/pkg/util"
)

type PaymentKind string

const (
	PaymentTypeOneTime      PaymentKind = "one_time"
	PaymentTypeSubscription PaymentKind = "subscription"
	PaymentTypeUnknown      PaymentKind = "unknown"
)

type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusPaid     PaymentStatus = "paid"
	PaymentStatusFailed   PaymentStatus = "failed"
	PaymentStatusRefunded PaymentStatus = "refunded"
)

type Payment struct {
	ID               string
	BillingAccountID string
	Amount           money.Money
	Status           PaymentStatus
	Kind             PaymentKind
	CreatedAt        time.Time
	PaidAt           *time.Time
}

type OneTimePaymentMetadata struct {
	ListingID string
	UserID    string
}

type SubscriptionPaymentType string

const (
	SubscriptionPaymentInstanceHosting SubscriptionPaymentType = "instance_hosting"
	SubscriptionPaymentUnknown         SubscriptionPaymentType = "unknown"
)

type SubscriptionPaymentMetadata struct {
	InstanceID       string
	CurrentPeriodEnd time.Time
	Type             SubscriptionPaymentType
}

func NewPendingPayment(
	billingAccountID string,
	amount money.Money,
) *Payment {

	return &Payment{
		ID:               util.GenerateUUID(),
		BillingAccountID: billingAccountID,
		Amount:           amount,
		Status:           PaymentStatusPending,
		Kind:             PaymentTypeUnknown,
		CreatedAt:        time.Now(),
	}
}
func (payment *Payment) MarkPaid(paidAt time.Time) {
	payment.Status = PaymentStatusPaid
	payment.PaidAt = &paidAt
}

func (payment *Payment) MarkFailed() {
	payment.Status = PaymentStatusFailed
}

func (payment *Payment) MarkRefunded() {
	payment.Status = PaymentStatusRefunded
}

type PaymentRepository interface {
	Save(payment *Payment) error
	GetByID(id string) (*Payment, error)
	ListByBillingAccount(accountID string) ([]*Payment, error)

	AttachOneTimeMetadata(paymentID string, metadata OneTimePaymentMetadata) error
	AttachSubscriptionMetadata(paymentID string, metadata SubscriptionPaymentMetadata) error

	GetOneTimeMetadata(paymentID string) (OneTimePaymentMetadata, error)
	GetSubscriptionMetadata(paymentID string) (SubscriptionPaymentMetadata, error)
}

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

	if payment.Kind != PaymentTypeUnknown {
		return fmt.Errorf("cannot add one-time metadata: payment kind is already %s", payment.Kind)
	}

	payment.Kind = PaymentTypeOneTime

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

	if payment.Kind != PaymentTypeUnknown {
		return fmt.Errorf("cannot add subscription metadata: payment kind is already %s", payment.Kind)
	}

	payment.Kind = PaymentTypeSubscription

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

	if payment.Kind != PaymentTypeOneTime {
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

	if payment.Kind != PaymentTypeSubscription {
		return SubscriptionPaymentMetadata{}, fmt.Errorf("payment does not have subscription metadata")
	}

	metadata, err := s.repo.GetSubscriptionMetadata(paymentID)
	if err != nil {
		return SubscriptionPaymentMetadata{}, err
	}

	return metadata, nil
}

func (s *PaymentService) ListByBillingAccount(accountID string) ([]*Payment, error) {
	return s.repo.ListByBillingAccount(accountID)
}

func (s *PaymentService) ListAllPaymentsByBillingAccount(accountID string) ([]*Payment, error) {
	return s.repo.ListByBillingAccount(accountID)
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

func (s *PaymentService) GetPaymentView(paymentID string) (*Payment, error) {
	return s.repo.GetByID(paymentID)
}
