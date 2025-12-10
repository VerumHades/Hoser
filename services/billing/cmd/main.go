package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"example.com/project/common/configuration"
	"example.com/project/handlers"
	"example.com/project/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config, err := configuration.Load[configuration.Configuration]()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	e := echo.New()

	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: strings.Join(config.AllowedOrigins, ","),
		AllowMethods: []string{http.MethodGet, http.MethodPost},
	}))

	// Initialize services
	billingService := services.NewBillingService(config.InterserverSecret, config.JWTSecret)

	// Routes
	e.POST("/bill", handlers.BillHandler(billingService))             // InterServer auth
	e.POST("/account", handlers.CreateAccountHandler(billingService)) // JWT auth
	e.GET("/payments", handlers.ListPaymentsHandler(billingService))  // JWT auth

	address := fmt.Sprintf("%s:%s", config.Address, config.Port)
	log.Printf("Starting server on %s", address)
	if err := e.Start(address); err != nil {
		log.Fatal(err)
	}
}
