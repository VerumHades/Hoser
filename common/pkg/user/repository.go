package user

// UserRepository defines persistence operations for users.
type UserRepository interface {
	Save(user *User) error
	GetByID(id string) (*User, error)
	GetByUsername(username string) (*User, error)
	Delete(id string) error
}
