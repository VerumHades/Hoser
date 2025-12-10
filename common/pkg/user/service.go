package user

import "common/pkg/util"

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
