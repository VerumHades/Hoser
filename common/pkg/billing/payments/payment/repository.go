package payment

type PaymentRepository interface {
	Save(payment *Payment) error
	GetByID(id string) (*Payment, error)
	ListByBillingAccount(accountID string) ([]*Payment, error)
}
