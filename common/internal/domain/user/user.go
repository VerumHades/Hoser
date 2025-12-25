package user

import (
	"common/internal/shared"
	"errors"
	"fmt"
)

// UserID is a value object representing a unique user identifier.
type UserID string

// User represents a system account internally as an aggregate root.
type User struct {
	id           UserID
	username     string
	passwordHash string
	developer    bool
}

// NewUser generates a new User entity with a generated ID.
func NewUser(username, passwordHash string, developer bool) (*User, error) {
	id := UserID(shared.GenerateUUID())
	return NewUserWithID(id, username, passwordHash, developer)
}

// NewUserWithID creates a User entity with an arbitrary ID.
// Intended for repository hydration or reconstruction.
func NewUserWithID(id UserID, username, passwordHash string, developer bool) (*User, error) {
	user := &User{
		id:           id,
		username:     username,
		passwordHash: passwordHash,
		developer:    developer,
	}
	if err := user.Validate(); err != nil {
		return nil, err
	}
	return user, nil
}

// ID returns the user's unique identifier.
func (u *User) ID() UserID {
	return u.id
}

// Username returns the user's username.
func (u *User) Username() string {
	return u.username
}

// Username returns the user's username.
func (u *User) PasswordHash() string {
	return u.passwordHash
}

// IsDeveloper returns whether the user has developer capabilities.
func (u *User) IsDeveloper() bool {
	return u.developer
}

// PromoteToDeveloper grants developer capabilities to the user.
func (u *User) PromoteToDeveloper() {
	u.developer = true
}

// DemoteFromDeveloper revokes developer capabilities from the user.
func (u *User) DemoteFromDeveloper() {
	u.developer = false
}

// ChangePassword updates the user's password hash.
func (u *User) ChangePassword(newHash string) error {
	if newHash == "" {
		return errors.New("password hash cannot be empty")
	}
	u.passwordHash = newHash
	return nil
}

// Validate ensures the user entity is in a valid state.
func (u *User) Validate() error {
	if u.id == "" {
		return errors.New("user ID cannot be empty")
	}
	if u.username == "" {
		return errors.New("username cannot be empty")
	}
	if u.passwordHash == "" {
		return errors.New("password hash cannot be empty")
	}
	return nil
}

// String implements a readable representation of the user.
func (u *User) String() string {
	return fmt.Sprintf("User[ID=%s, Username=%s, Developer=%v]", u.id, u.username, u.developer)
}
