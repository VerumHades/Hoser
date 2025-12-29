package user

import (
	"common/internal/shared"
	"context"
)

// UserCommandRepository defines write operations for users.
type UserCommandRepository interface {
	// Create persists a new user within a transaction.
	Create(
		ctx context.Context,
		transaction shared.Transaction,
		user *User,
	) (*User, error)

	// Update modifies an existing user within a transaction.
	Update(
		ctx context.Context,
		transaction shared.Transaction,
		user *User,
	) (*User, error)

	// Delete removes a user by its ID within a transaction.
	Delete(
		ctx context.Context,
		transaction shared.Transaction,
		userID shared.UserID,
	) error
}

// UserQueryRepository defines read-only operations for users.
type UserQueryRepository interface {
	// GetByID retrieves a user by its ID.
	GetByID(
		ctx context.Context,
		userID shared.UserID,
	) (*User, error)

	Exists(
		ctx context.Context,
		userID shared.UserID,
	) (bool, error)

	// GetByUsername retrieves a user by their username.
	GetByUsername(
		ctx context.Context,
		username string,
	) (*User, error)
}
