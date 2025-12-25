package billing

import (
	"common/internal/shared"
	"time"
)

type BillingAccountRepository interface {
	Save(account *BillingAccount) error

	GetByID(accountID shared.BillingAccountID) (*BillingAccount, error)
	ListByOwner(ownerID shared.UserID) ([]*BillingAccount, error)
}

type PaymentRepository interface {
	Save(payment *Payment) error
	GetByID(id shared.PaymentID) (*Payment, error)
	ListByBillingAccount(accountID shared.BillingAccountID) ([]*Payment, error)
}

// HardwareCostRepository defines persistence operations for hardware cost rates.
type HardwareCostRepository interface {
	Save(rate *HardwareCostRate) error
	GetActiveRate(at time.Time) (*HardwareCostRate, error)
	ListAll() ([]*HardwareCostRate, error)
}
