package user

import "common/internal/shared"

// UserRepository defines persistence operations for users.
type UserRepository interface {
	Save(user *User) (*User, error)

	GetByID(id shared.UserID) (*User, error)
	GetByUsername(username string) (*User, error)
	Delete(id shared.UserID) error
}
