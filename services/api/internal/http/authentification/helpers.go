package authentification

import (
	"common/pkg/shared"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandlerFunc func(userID shared.UserID, context echo.Context) error

func WithAuthenticatedUser(userHandler UserHandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, ok := c.Get("user_id").(string)

		if !ok || userID == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing authenticated user")
		}
		return userHandler(shared.UserID(userID), c)
	}
}

func GetAuthenticatedUserID(context echo.Context) (shared.UserID, error) {
	userID, ok := context.Get("user_id").(shared.UserID)
	if !ok || userID == "" {
		return "", echo.NewHTTPError(http.StatusUnauthorized)
	}

	return userID, nil
}
