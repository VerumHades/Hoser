package domain

import (
	"errors"
	"time"

	"common/internal/shared"
)

// BillingAccount represents a user's billing account.
type BillingAccount struct {
	id                shared.BillingAccountID
	ownerID           shared.UserID
	status            BillingAccountStatus
	paymentProvider   PaymentProvider
	providerAccountID string
	createdAt         time.Time
}

// NewBillingAccount creates a new billing account with a generated ID and the Active status.
func NewBillingAccount(ownerID shared.UserID, paymentProvider PaymentProvider, providerAccountID string) (*BillingAccount, error) {
	id := shared.BillingAccountID(shared.GenerateUUID())
	return NewBillingAccountWithID(id, ownerID, Active, paymentProvider, providerAccountID, time.Now().UTC())
}

// NewBillingAccountWithID creates a billing account with an existing ID (for repository hydration).
func NewBillingAccountWithID(id shared.BillingAccountID, ownerID shared.UserID, status BillingAccountStatus, paymentProvider PaymentProvider, providerAccountID string, createdAt time.Time) (*BillingAccount, error) {
	account := &BillingAccount{
		id:                id,
		ownerID:           ownerID,
		status:            status,
		paymentProvider:   paymentProvider,
		providerAccountID: providerAccountID,
		createdAt:         createdAt,
	}
	if err := account.Validate(); err != nil {
		return nil, err
	}
	return account, nil
}

// ID returns the unique identifier.
func (b *BillingAccount) ID() shared.BillingAccountID {
	return b.id
}

// OwnerID returns the owner's user ID.
func (b *BillingAccount) OwnerID() shared.UserID {
	return b.ownerID
}

// Status returns the current account status.
func (b *BillingAccount) Status() BillingAccountStatus {
	return b.status
}

func (b *BillingAccount) IsActive() bool {
	return b.status == Active
}

// PaymentProvider returns the associated payment provider.
func (b *BillingAccount) PaymentProvider() PaymentProvider {
	return b.paymentProvider
}

// ProviderAccountID returns the external payment provider account ID.
func (b *BillingAccount) ProviderAccountID() string {
	return b.providerAccountID
}

// CreatedAt returns the creation timestamp.
func (b *BillingAccount) CreatedAt() time.Time {
	return b.createdAt
}

// Activate sets the account status to Active.
func (b *BillingAccount) Activate() {
	b.status = Active
}

// Suspend sets the account status to Suspended.
func (b *BillingAccount) Suspend() {
	b.status = Suspended
}

// Close sets the account status to Closed.
func (b *BillingAccount) Close() {
	b.status = Closed
}

// Validate ensures the billing account is in a valid state.
func (b *BillingAccount) Validate() error {
	if b.id == "" {
		return errors.New("billing account ID cannot be empty")
	}
	if b.ownerID == "" {
		return errors.New("owner ID cannot be empty")
	}
	if b.paymentProvider == "" {
		return errors.New("payment provider cannot be nil")
	}
	if b.providerAccountID == "" {
		return errors.New("provider account ID cannot be empty")
	}
	if b.createdAt.IsZero() {
		return errors.New("createdAt cannot be zero")
	}
	switch b.status {
	case Active, Suspended, Closed:
		// valid
	default:
		return errors.New("invalid billing account status")
	}
	return nil
}
