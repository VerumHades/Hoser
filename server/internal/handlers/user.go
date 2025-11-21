package handlers

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest defines the expected login payload
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse defines the response with JWT
type LoginResponse struct {
	Token string `json:"token"`
}

func (app *App) LoginHandler(c echo.Context) error {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	user, err := app.DatabaseInteractor.GetUserByName(req.Username)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash()), []byte(req.Password)) != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	// Create JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID(),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(app.JWTSecret))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate token")
	}

	// Set cookie
	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = tokenString
	cookie.HttpOnly = true
	cookie.Secure = true // HTTPS only
	cookie.Path = "/"
	cookie.SameSite = http.SameSiteNoneMode // allow cross-origin
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Logged in successfully",
	})
}

func (app *App) LogoutHandler(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = ""
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.Path = "/"
	cookie.Expires = time.Unix(0, 0)
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}

// UserDataHandler returns info about the currently authenticated user
func (app *App) UserDataHandler(c echo.Context) error {
	user, err := app.GetUserFromContext(c)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, struct {
		Username    string `json:"username"`
		IsDeveloper bool   `json:"isDeveloper"`
	}{
		Username:    user.Username(),
		IsDeveloper: user.IsDeveloper(),
	})
}
