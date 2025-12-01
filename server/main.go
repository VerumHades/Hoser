package main

import (
	"fmt"
	"os"
	"path/filepath"
	"server/internal/configuration"
	"server/internal/database"
	"server/internal/handlers"

	"github.com/labstack/echo/v4"
)

func main() {
	// ---------------------------
	// Load configuration
	// ---------------------------
	runningConfiguration := configuration.Load()

	app := &handlers.App{
		RunningConfiguration: &runningConfiguration,
		DatabaseInteractor:   &database.DummyStore{},
		JWTSecret:            []byte(runningConfiguration.SessionSecret),
	}

	e := echo.New()
	// ---------------------------
	// Middleware
	// ---------------------------
	e.Use(app.CORSMiddleware())
	//e.Use(middleware.Logger())
	//e.Use(middleware.Recover())

	// ---------------------------
	// Public routes
	// ---------------------------
	public := e.Group("")
	public.GET("/rentals/public", app.PublicListingsHandler)
	public.POST("/login", app.LoginHandler)
	public.POST("/logout", app.LogoutHandler)

	// ---------------------------
	// Authenticated routes
	// ---------------------------
	auth := e.Group("/user")

	auth.GET("/data", app.UserDataHandler)
	auth.GET("/rentals", app.UserRentalsHandler)
	auth.POST("/rent", app.UserRentHandler)

	// ---------------------------
	// Developer-only routes
	// ---------------------------
	dev := e.Group("/developer/listing")
	dev.Use(app.DeveloperOnlyMiddleware)

	dev.GET("", app.DeveloperListingsHandler)
	dev.POST("", app.DeveloperAddListingHandler)
	dev.PUT("", app.DeveloperAlterListingHandler)
	dev.DELETE("", app.DeveloperDeleteListingHandler)

	price := dev.Group("/price")
	price.Use(app.DeveloperOnlyMiddleware)
	price.POST("", app.DeveloperAddListingPricingHandler)
	price.PUT("", app.DeveloperUpdateListingPricingHandler)
	price.DELETE("", app.DeveloperDeleteListingPriceHandler)

	// ---------------------------
	// Serve React client
	// ---------------------------
	fsRoot := runningConfiguration.ClientDirectory
	e.GET("/*", func(c echo.Context) error {
		path := filepath.Join(fsRoot, c.Request().URL.Path)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// Fallback to index.html
			return c.File(filepath.Join(fsRoot, "index.html"))
		}
		return c.File(path)
	})

	// ---------------------------
	// Start server
	// ---------------------------
	address := fmt.Sprintf("%s:%s", runningConfiguration.Address, runningConfiguration.Port)
	fmt.Printf("Server running on http://%s\n", address)
	e.Logger.Fatal(e.Start(address))
}
