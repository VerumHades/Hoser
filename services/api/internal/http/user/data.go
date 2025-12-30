package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// UserDataHandler returns info about the currently authenticated user
func (api *UserAPI) UserDataHandler(context echo.Context) error {
	userID, _ := GetAuthenticatedUserID(context)
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
