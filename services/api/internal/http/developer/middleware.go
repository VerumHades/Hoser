package developer

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// DeveloperCheckMiddleware ensures the authenticated user is a developer.
func (api *DeveloperListingAPI) DeveloperCheckMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		// Retrieve authenticated user ID from context
		userIDValue := context.Get("user_id")
		userID, ok := userIDValue.(string)
		if !ok || userID == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing authenticated user")
		}

		// Check if the user has developer privileges
		isDeveloper, err := api.UserService.IsUserDeveloper(userID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to verify user role")
		}
		if !isDeveloper {
			return echo.NewHTTPError(http.StatusForbidden, "User is not a developer")
		}

		// User is a developer, proceed to next handler
		return next(context)
	}
}
