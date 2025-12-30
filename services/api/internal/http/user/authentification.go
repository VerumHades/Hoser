package user

import (
	"common/pkg/shared"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

// AuthenticationMiddleware validates the JWT and injects the authenticated user ID into the request context.
func (api *UserAPI) AuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		cookie, cookieError := context.Cookie("jwt")
		if cookieError != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing token")
		}

		token, parseError := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			return []byte(api.runningConfiguration.JWTSecret), nil
		})
		if parseError != nil || !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		claims, claimsOk := token.Claims.(jwt.MapClaims)
		if !claimsOk {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		userID, userIDOk := claims["user_id"].(string)
		if !userIDOk || userID == "" {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		context.Set("user_id", userID)

		return next(context)
	}
}

func GetAuthenticatedUserID(context echo.Context) (shared.UserID, error) {
	userID, ok := context.Get("user_id").(shared.UserID)
	if !ok || userID == "" {
		return "", echo.NewHTTPError(http.StatusUnauthorized)
	}

	return userID, nil
}
