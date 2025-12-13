package account

type BillingAccountRepository interface {
	Save(account *BillingAccount) error

	GetByID(accountID string) (*BillingAccount, error)
	ListByOwner(ownerID string) ([]*BillingAccount, error)
}
