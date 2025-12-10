package main

import (
	"bytes"
	"common/pkg/configuration"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"engine/internal/reconciliation"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// ControllerHandler handles HTTP routes and interactions with the controller
type ControllerHandler struct {
	controller    reconciliation.Controller
	configuration configuration.Configuration
}

// GetUserIDFromContext extracts the user ID from JWT stored in cookies
func (handler *ControllerHandler) GetUserIDFromContext(c echo.Context) (string, error) {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Missing token")
	}

	token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
		return []byte(handler.configuration.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", echo.NewHTTPError(http.StatusUnauthorized)
	}

	return claims["user_id"].(string), nil
}

// computeHmacSignature generates HMAC-SHA256 signature of payload with shared secret
func computeHmacSignature(sharedSecret string, payloadToSign string) string {
	mac := hmac.New(sha256.New, []byte(sharedSecret))
	mac.Write([]byte(payloadToSign))
	return hex.EncodeToString(mac.Sum(nil))
}

// InterserverSignatureMiddleware validates HMAC signatures for internal requests
func (handler *ControllerHandler) InterserverSignatureMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			signatureHeader := context.Request().Header.Get("X-Signature")
			timestampHeader := context.Request().Header.Get("X-Timestamp")

			if signatureHeader == "" || timestampHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing signature headers")
			}

			timestampValue, err := strconv.ParseInt(timestampHeader, 10, 64)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid timestamp")
			}

			if time.Since(time.Unix(timestampValue, 0)) > 5*time.Minute {
				return echo.NewHTTPError(http.StatusUnauthorized, "request expired")
			}

			requestBodyBytes, err := io.ReadAll(context.Request().Body)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
			}
			context.Request().Body = io.NopCloser(bytes.NewBuffer(requestBodyBytes))

			stringToSign := strings.Join([]string{
				context.Request().Method,
				context.Request().URL.Path,
				timestampHeader,
				string(requestBodyBytes),
			}, "\n")

			expectedSignature := computeHmacSignature(handler.configuration.InterserverSecret, stringToSign)

			if !hmac.Equal([]byte(signatureHeader), []byte(expectedSignature)) {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid signature")
			}

			return next(context)
		}
	}
}

// RegisterRoutes sets up public and internal HTTP routes
func (handler *ControllerHandler) RegisterRoutes(echoServer *echo.Echo) {
	publicGroup := echoServer.Group("/instances")
	internalGroup := echoServer.Group("/instances", handler.InterserverSignatureMiddleware())

	publicGroup.GET("", handler.ListInstances)
	publicGroup.GET("/:identifier", handler.GetInstance)
	publicGroup.PUT("/:identifier", handler.UpdateInstance)

	internalGroup.POST("", handler.CreateInstance)
	internalGroup.DELETE("/:identifier", handler.DeleteInstance)
	internalGroup.POST("/:identifier/lifetime", handler.UpdateInstanceLifetime) // NEW ENDPOINT
}

/*
	Internal endpoints
*/

// CreateInstance handles internal POST requests to create an instance
func (handler *ControllerHandler) CreateInstance(context echo.Context) error {
	var instanceSpec reconciliation.InstanceCreationSpecification

	if err := context.Bind(&instanceSpec); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	identifier, err := handler.controller.CreateInstance(instanceSpec)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusCreated, map[string]string{"id": identifier})
}

// DeleteInstance handles internal DELETE requests to remove an instance
func (handler *ControllerHandler) DeleteInstance(context echo.Context) error {
	instanceIdentifier := context.Param("identifier")

	if err := handler.controller.DeleteInstance(instanceIdentifier); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return context.NoContent(http.StatusNoContent)
}

// UpdateInstanceLifetime handles inter-server POST requests to update expiration
func (handler *ControllerHandler) UpdateInstanceLifetime(context echo.Context) error {
	instanceIdentifier := context.Param("identifier")

	var request struct {
		NewExpiration int `json:"new_expiration"`
	}

	if err := context.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if request.NewExpiration <= int(time.Now().Unix()) {
		return echo.NewHTTPError(http.StatusBadRequest, "new expiration must be in the future")
	}

	if err := handler.controller.ExtendInstanceLifetime(instanceIdentifier, request.NewExpiration); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, struct {
		identifier     string
		new_expiration int
	}{
		identifier:     instanceIdentifier,
		new_expiration: request.NewExpiration,
	})
}

/*
	Public endpoints
*/

// ListInstances lists all instances for the authenticated user
func (handler *ControllerHandler) ListInstances(context echo.Context) error {
	userIdentifier, err := handler.GetUserIDFromContext(context)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user id")
	}

	instances, err := handler.controller.ListInstancesByOwner(userIdentifier)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, instances)
}

// GetInstance retrieves a single instance for the authenticated user
func (handler *ControllerHandler) GetInstance(context echo.Context) error {
	userIdentifier, err := handler.GetUserIDFromContext(context)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user id")
	}

	instanceIdentifier := context.Param("identifier")
	instanceState, err := handler.controller.GetInstanceForOwner(instanceIdentifier, userIdentifier)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return context.JSON(http.StatusOK, instanceState)
}

// UpdateInstance updates a user's instance based on provided specification
func (handler *ControllerHandler) UpdateInstance(context echo.Context) error {
	userIdentifier, err := handler.GetUserIDFromContext(context)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user id")
	}

	instanceIdentifier := context.Param("identifier")
	var updateSpec reconciliation.InstanceUpdateSpecification

	if err := context.Bind(&updateSpec); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := handler.controller.UpdateInstanceForOwner(instanceIdentifier, userIdentifier, updateSpec); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	updatedState, err := handler.controller.GetInstanceForOwner(instanceIdentifier, userIdentifier)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, updatedState)
}

func main() {
	config, err := configuration.Load[configuration.Configuration]()
	if err != nil {
		log.Fatal(err)
	}

	echoServer := echo.New()

	if len(config.AllowedOrigins) > 0 && config.AllowedOrigins[0] != "" {
		echoServer.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: config.AllowedOrigins,
			AllowMethods: []string{echo.GET, echo.HEAD, echo.OPTIONS},
		}))
	}

	var controller reconciliation.Controller = reconciliation.DockerController{}
	handler := ControllerHandler{controller: controller, configuration: config}
	handler.RegisterRoutes(echoServer)

	address := fmt.Sprintf("%s:%s", config.Address, config.Port)
	fmt.Printf("Frontend server running on http://%s\n", address)
	echoServer.Logger.Fatal(echoServer.Start(address))
}
