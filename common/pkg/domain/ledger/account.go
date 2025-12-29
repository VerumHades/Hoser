package ledger

import (
	"common/pkg/shared"
	"fmt"
	"time"
)

// AccountType defines the role of an account in the ledger.
type AccountType string

const (
	AccountTypeUserLiability   AccountType = "user_liability"
	AccountTypePlatformRevenue AccountType = "platform_revenue"
	AccountTypeClearing        AccountType = "clearing"
	AccountTypeSuspense        AccountType = "suspense"
)

type AccountOwnerType string

const (
	AccountOwnerPlatform  AccountOwnerType = "platform"
	AccountOwnerUser      AccountOwnerType = "user"
	AccountOwnerProcessor AccountOwnerType = "processor"
)

// Account is a domain entity representing a ledger account.
// Fields are private; all creation and modification enforce invariants.
type Account struct {
	id          shared.AccountID
	accountType AccountType
	ownerType   AccountOwnerType
	ownerID     string
	createdAt   time.Time
}

// NewAccount constructs a new Account while enforcing invariants.
func NewAccount(id shared.AccountID, accountType AccountType, ownerType AccountOwnerType, ownerID string) (*Account, error) {
	if id == "" {
		return nil, fmt.Errorf("account ID cannot be empty")
	}
	if ownerType != AccountOwnerPlatform && ownerType != AccountOwnerUser && ownerType != AccountOwnerProcessor {
		return nil, fmt.Errorf("invalid owner type: %s", ownerType)
	}
	if ownerID == "" && ownerType != AccountOwnerPlatform {
		return nil, fmt.Errorf("ownerID cannot be empty for non-platform accounts")
	}
	switch accountType {
	case AccountTypeUserLiability, AccountTypePlatformRevenue, AccountTypeClearing, AccountTypeSuspense:
		// valid
	default:
		return nil, fmt.Errorf("invalid account type: %s", accountType)
	}

	return &Account{
		id:          id,
		accountType: accountType,
		ownerType:   ownerType,
		ownerID:     ownerID,
		createdAt:   time.Now(),
	}, nil
}

// Getters for read access
func (a *Account) ID() shared.AccountID {
	return a.id
}

func (a *Account) Type() AccountType {
	return a.accountType
}

func (a *Account) OwnerType() AccountOwnerType {
	return a.ownerType
}

func (a *Account) OwnerID() string {
	return a.ownerID
}

func (a *Account) CreatedAt() time.Time {
	return a.createdAt
}

// Example of a domain method that enforces invariants
func (a *Account) ChangeOwner(newOwnerType AccountOwnerType, newOwnerID string) error {
	if newOwnerType != AccountOwnerUser && newOwnerType != AccountOwnerPlatform && newOwnerType != AccountOwnerProcessor {
		return fmt.Errorf("invalid new owner type: %s", newOwnerType)
	}
	if newOwnerID == "" && newOwnerType != AccountOwnerPlatform {
		return fmt.Errorf("newOwnerID cannot be empty for non-platform accounts")
	}
	a.ownerType = newOwnerType
	a.ownerID = newOwnerID
	return nil
}
