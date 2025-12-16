package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// =================== MIDDLEWARE ===================

func (app *App) DeveloperOnlyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := app.GetUserIDFromContext(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		isDeveloper, err := app.UserAppService.IsUserDeveloper(userID)

		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		if !isDeveloper {
			return echo.NewHTTPError(http.StatusUnauthorized, "Developer access required")
		}
		return next(c)
	}
}

func (app *App) CORSMiddleware() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     app.RunningConfiguration.AllowedOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	})
}
