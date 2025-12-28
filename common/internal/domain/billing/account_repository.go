package billing

import (
	"common/internal/shared"
	"context"
)

type BillingAccountCommandRepository interface {
	Create(
		context context.Context,
		transaction shared.Transaction,
		account *BillingAccount,
	) error

	Update(
		context context.Context,
		transaction shared.Transaction,
		account *BillingAccount,
	) error
}

type BillingAccountQueryRepository interface {
	GetByID(
		context context.Context,
		accountID shared.BillingAccountID,
	) (account *BillingAccount, err error)

	FetchNextBatchByOwner(
		context context.Context,
		ownerID shared.UserID,
		request shared.BatchRequest,
	) (accounts []*BillingAccount, nextCursor shared.Cursor, err error)
}
