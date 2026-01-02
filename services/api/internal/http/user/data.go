package userapi

import (
	"api/internal/http/authentification"
	"common/pkg/domain/user"
	"common/pkg/shared"
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type APIUserDataQueryService interface {
	GetByID(context context.Context, userID shared.UserID) (*user.User, error)
}

type UserProfileAPI struct {
	userQueryService APIUserDataQueryService
}

// NewUserProfileAPI constructs a UserProfileAPI.
func NewUserProfileAPI(
	userQueryService APIUserDataQueryService,
) *UserProfileAPI {
	return &UserProfileAPI{
		userQueryService: userQueryService,
	}
}

func (api *UserProfileAPI) RegisterRoutes(group *echo.Group) {
	group.GET("/data", authentification.WithAuthenticatedUser(api.UserDataHandler))
}

// UserDataHandler returns info about the currently authenticated user
func (api *UserProfileAPI) UserDataHandler(userID shared.UserID, context echo.Context) error {
	requestContext := context.Request().Context()

	user, err := api.userQueryService.GetByID(requestContext, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	return context.JSON(http.StatusOK, struct {
		Username    string `json:"username"`
		IsDeveloper bool   `json:"isDeveloper"`
	}{
		Username:    user.Username(),
		IsDeveloper: user.IsDeveloper(),
	})
}
