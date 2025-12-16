package user

import "common/pkg/util"

// User represents a system account internally.
type User struct {
	id           string
	username     string
	passwordHash string
	developer    bool
}

func NewUser(id, username, passwordHash string, developer bool) *User {
	return &User{
		id:           id,
		username:     username,
		passwordHash: passwordHash,
		developer:    developer,
	}
}

// UserPublicView exposes safe read-only public information about a user.
type UserPublicView struct {
	ID       string
	Username string
	IsDev    bool
}

// UserPrivateView exposes sensitive internal information about a user.
type UserPrivateView struct {
	ID           string
	Username     string
	IsDev        bool
	PasswordHash string
}

// ToPublicView converts a User into a UserPublicView.
func (entity *User) ToPublicView() *UserPublicView {
	publicView := &UserPublicView{
		ID:       entity.id,
		Username: entity.username,
		IsDev:    entity.developer,
	}
	return publicView
}

// ToPrivateView converts a User into a UserPrivateView.
func (entity *User) ToPrivateView() *UserPrivateView {
	privateView := &UserPrivateView{
		ID:           entity.id,
		Username:     entity.username,
		IsDev:        entity.developer,
		PasswordHash: entity.passwordHash,
	}
	return privateView
}

// UserRepository defines persistence operations for users.
type UserRepository interface {
	Save(user *User) error
	GetByID(id string) (*User, error)
	GetByUsername(username string) (*User, error)
	Delete(id string) error
}

// UserService handles business logic strictly related to users.
type UserService struct {
	repo UserRepository
}

// NewUserService creates a new instance.
func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

// CreateUser creates a new user account.
func (s *UserService) CreateUser(username, passwordHash string, developer bool) (*User, error) {
	user := &User{
		id:           util.GenerateUUID(),
		username:     username,
		passwordHash: passwordHash,
		developer:    developer,
	}
	err := s.repo.Save(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// CheckDeveloper returns whether the user is a developer.
func (s *UserService) CheckDeveloper(userID string) (bool, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return false, err
	}
	return user.developer, nil
}

// CheckDeveloper returns whether the user is a developer.
func (s *UserService) GetUser(userID string) (*User, error) {
	return s.repo.GetByID(userID)
}
