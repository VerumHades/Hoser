package auth

import (
	"errors"

	"common/pkg/user"

	"golang.org/x/crypto/bcrypt"
)

type AuthenticationService struct {
	userRepository user.UserRepository
}

func NewAuthenticationService(userRepository user.UserRepository) *AuthenticationService {
	return &AuthenticationService{
		userRepository: userRepository,
	}
}

// AuthenticateUser verifies a user's credentials and returns the user if valid.
func (service *AuthenticationService) AuthenticateUser(username string, password string) (*user.User, error) {
	fetchedUser, fetchError := service.userRepository.GetByUsername(username)
	if fetchError != nil {
		return nil, errors.New("invalid credentials")
	}

	passwordError := bcrypt.CompareHashAndPassword([]byte(fetchedUser.ToPrivateView().PasswordHash), []byte(password))
	if passwordError != nil {
		return nil, errors.New("invalid credentials")
	}

	return fetchedUser, nil
}
