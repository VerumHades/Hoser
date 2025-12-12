package account

type BillingAccountRepository interface {
	Save(account *BillingAccount) error
	GetByID(accountID string) (*BillingAccount, error)
	GetByOwnerID(ownerID string) (*BillingAccount, error)
}
