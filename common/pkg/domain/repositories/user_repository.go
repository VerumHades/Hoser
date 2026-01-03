package repositories

import (
	"common/pkg/domain/entities/user"
	"common/pkg/shared"
	"context"
)

// UserCommandRepository defines write operations for users.
type UserCommandRepository interface {
	// Create persists a new user within a transaction.
	Create(
		ctx context.Context,

		user *user.User,
	) (*user.User, error)

	// Update modifies an existing user within a transaction.
	Update(
		ctx context.Context,

		user *user.User,
	) (*user.User, error)

	// Delete removes a user by its ID within a transaction.
	Delete(
		ctx context.Context,

		userID shared.UserID,
	) error
}

// UserQueryRepository defines read-only operations for users.
type UserQueryRepository interface {
	// GetByID retrieves a user by its ID.
	GetByID(
		ctx context.Context,
		userID shared.UserID,
	) (*user.User, error)

	Exists(
		ctx context.Context,
		userID shared.UserID,
	) (bool, error)

	// GetByUsername retrieves a user by their username.
	GetByUsername(
		ctx context.Context,
		username string,
	) (*user.User, error)
}
