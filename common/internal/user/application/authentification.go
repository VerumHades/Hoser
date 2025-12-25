package application

import (
	"common/internal/user/domain"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// AuthenticationService handles user authentication workflows.
type AuthenticationService struct {
	userRepository domain.UserRepository
}

// NewAuthenticationService creates a new AuthenticationService.
func NewAuthenticationService(userRepository domain.UserRepository) *AuthenticationService {
	return &AuthenticationService{
		userRepository: userRepository,
	}
}

// AuthenticateUser verifies a user's credentials and returns the user if valid.
func (service *AuthenticationService) AuthenticateUser(username string, password string) (*domain.User, error) {
	user, err := service.userRepository.GetByUsername(username)
	if err != nil || user == nil {
		return nil, errors.New("invalid credentials")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash()), []byte(password)) != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
