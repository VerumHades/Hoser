package service

type BillingService interface {
	CreateAccount(userID string) error
	Bill(listingID string, billingID string) error
	ListPayments(userID string) ([]string, error)
}
