package auth

import (
	"common/pkg/domain/entities/user"
	"common/pkg/domain/repositories"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// AuthenticationService handles user authentication workflows.
type AuthenticationService struct {
	userQueryRepository repositories.UserQueryRepository
}

// NewAuthenticationService creates a new AuthenticationService.
func NewAuthenticationService(userRepository repositories.UserQueryRepository) *AuthenticationService {
	return &AuthenticationService{
		userQueryRepository: userRepository,
	}
}

// AuthenticateUser verifies a user's credentials and returns the user if valid.
func (service *AuthenticationService) AuthenticateUser(context context.Context, username string, password string) (*user.User, error) {
	user, err := service.userQueryRepository.GetByUsername(context, username)
	if err != nil || user == nil {
		return nil, errors.New("invalid credentials")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash()), []byte(password)) != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
