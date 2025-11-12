package handlers

import (
	"net/http"
	"server/internal/configuration"
	"server/internal/database"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type App struct {
	DatabaseInteractor   database.Interactor
	RunningConfiguration *configuration.Configuration
	JWTSecret            []byte
}

// =================== APP INIT ===================

func NewApp(db database.Interactor, cfg *configuration.Configuration, jwtSecret []byte) *App {
	return &App{
		DatabaseInteractor:   db,
		RunningConfiguration: cfg,
		JWTSecret:            jwtSecret,
	}
}

// =================== HELPER FUNCTIONS ===================

// GetUserFromContext retrieves the authenticated user from Echo context (JWT claims)
func (app *App) GetUserFromContext(c echo.Context) (database.User, error) {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Missing token")
	}

	token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
		return []byte(app.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, echo.NewHTTPError(http.StatusUnauthorized)
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, echo.NewHTTPError(http.StatusUnauthorized)
	}

	user, err := app.DatabaseInteractor.GetUserByID(userID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized)
	}

	return user, nil
}

// =================== TEST APP ===================

func NewTestApp() *App {
	db := &database.DummyInteractor{}
	cfg := &configuration.Configuration{
		AllowedOrigins: []string{"http://localhost"},
	}
	jwtSecret := []byte("test-secret-key")

	return NewApp(db, cfg, jwtSecret)
}
