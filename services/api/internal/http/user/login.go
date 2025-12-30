package user

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
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

func (api *UserAPI) LoginHandler(context echo.Context) error {
	requestContext := context.Request().Context()

	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	req := new(LoginRequest)
	if err := context.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	user, err := api.userAuthentificationService.AuthenticateUser(requestContext, req.Username, req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	// Create JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID(),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(api.runningConfiguration.JWTSecret))
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
	context.SetCookie(cookie)

	return context.JSON(http.StatusOK, map[string]string{
		"message": "Logged in successfully",
	})
}

func (api *UserAPI) LogoutHandler(c echo.Context) error {
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
