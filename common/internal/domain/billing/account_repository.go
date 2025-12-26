package billing

import (
	"common/internal/shared"
)

type BillingAccountRepository interface {
	Save(account *BillingAccount) (*BillingAccount, error)

	GetByID(accountID shared.BillingAccountID) (*BillingAccount, error)

	FetchNextBatchByOwner(
		ownerID shared.UserID,
		lastSeenAccountID shared.BillingAccountID,
		maximumBatchSize int,
	) ([]*BillingAccount, error)
}
