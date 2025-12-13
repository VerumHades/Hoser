package payment

type PaymentRepository interface {
	Save(payment *Payment) error
	GetByID(id string) (*Payment, error)
	ListByBillingAccount(accountID string) ([]*Payment, error)

	AttachOneTimeMetadata(paymentID string, metadata OneTimePaymentMetadata) error
	AttachSubscriptionMetadata(paymentID string, metadata SubscriptionPaymentMetadata) error

	GetOneTimeMetadata(paymentID string) (OneTimePaymentMetadata, error)
	GetSubscriptionMetadata(paymentID string) (SubscriptionPaymentMetadata, error)
}
