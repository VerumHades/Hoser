package repositories

import (
	"common/pkg/domain/entities/accounting"
	"common/pkg/shared"
	"context"
)

// AccountCommandRepository defines write operations for accounts.
type AccountCommandRepository interface {
	// Create persists a new account within a transaction.
	Create(
		ctx context.Context,

		account *accounting.Account,
	) (*accounting.Account, error)

	// Update modifies an existing account within a transaction.
	Update(
		ctx context.Context,

		account *accounting.Account,
	) (*accounting.Account, error)

	// Delete removes an account by its ID within a transaction.
	Delete(
		ctx context.Context,

		accountID shared.AccountID,
	) error
}

// AccountQueryRepository defines read-only operations for accounts.
type AccountQueryRepository interface {
	// GetByID retrieves an account by its ID.
	GetByID(
		ctx context.Context,
		accountID shared.AccountID,
	) (*accounting.Account, error)

	// GetByOwner retrieves all accounts owned by a specific owner (user, platform, or processor).
	GetByOwner(
		ctx context.Context,
		ownerType accounting.AccountOwnerType,
		ownerID string,
	) ([]*accounting.Account, error)

	GetFirstByOwner(
		ctx context.Context,
		ownerType accounting.AccountOwnerType,
		ownerID string,
	) (*accounting.Account, error)

	// GetByType retrieves all accounts of a specific account type.
	GetByType(
		ctx context.Context,
		accountType accounting.AccountType,
	) ([]*accounting.Account, error)
}
