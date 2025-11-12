package handlers

import (
	"fmt"
	"net/http"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// =================== MIDDLEWARE ===================

func (app *App) JWTMiddleware() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey:  app.JWTSecret,
		ContextKey:  "user",
		TokenLookup: "header:Authorization",
	})
}

func (app *App) DeveloperOnlyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := app.GetUserFromContext(c)
		fmt.Println(user)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		if !user.IsDeveloper() {
			return echo.NewHTTPError(http.StatusUnauthorized, "Developer access required")
		}
		return next(c)
	}
}

func (app *App) CORSMiddleware() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     app.RunningConfiguration.AllowedOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderAuthorization, echo.HeaderContentType},
		AllowCredentials: true,
	})
}
