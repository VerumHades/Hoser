package auth

import (
	"errors"

	"common/pkg/user"

	"golang.org/x/crypto/bcrypt"
)

type AuthenticationService struct {
	userService *user.UserService
}

func NewAuthenticationService(userService *user.UserService) *AuthenticationService {
	return &AuthenticationService{
		userService: userService,
	}
}

// AuthenticateUser verifies a user's credentials and returns the user if valid.
func (service *AuthenticationService) AuthenticateUser(username string, password string) (*user.User, error) {
	fetchedUser, fetchError := service.userService.GetByUsername(username)
	if fetchError != nil {
		return nil, errors.New("invalid credentials")
	}

	passwordError := bcrypt.CompareHashAndPassword([]byte(fetchedUser.ToPrivateView().PasswordHash), []byte(password))
	if passwordError != nil {
		return nil, errors.New("invalid credentials")
	}

	return fetchedUser, nil
}
