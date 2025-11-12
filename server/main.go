package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"server/internal/configuration"
	"server/internal/database"
	"server/internal/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// ---------------------------
	// Load configuration
	// ---------------------------
	runningConfiguration := configuration.Load()

	app := &handlers.App{
		RunningConfiguration: &runningConfiguration,
		DatabaseInteractor:   &database.DummyInteractor{},
		JWTSecret:            []byte(runningConfiguration.SessionSecret),
	}

	e := echo.New()

	// ---------------------------
	// Middleware
	// ---------------------------
	e.Use(app.CORSMiddleware())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// ---------------------------
	// Public routes
	// ---------------------------
	public := e.Group("")
	public.GET("/rentals/public", func(c echo.Context) error {
		options := database.RentalQueryOptions{}
		q := c.QueryParam("q")
		if q != "" {
			options.Text = q
		}

		response, err := app.DatabaseInteractor.QueryPublicListings(&options)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.JSON(http.StatusOK, response)
	})
	public.POST("/login", app.LoginHandler)
	public.POST("/logout", app.LogoutHandler)

	// ---------------------------
	// Authenticated routes
	// ---------------------------
	auth := e.Group("")

	auth.GET("/user/data", app.UserDataHandler)
	auth.GET("/user/rentals", app.UserRentalsHandler)
	auth.POST("/user/rent", app.UserRentHandler)

	// ---------------------------
	// Developer-only routes
	// ---------------------------
	dev := e.Group("")
	dev.Use(app.DeveloperOnlyMiddleware)

	dev.GET("/developer/listings", app.DeveloperListingsHandler)
	dev.POST("/developer/listing", app.DeveloperAddListingHandler)
	dev.PUT("/developer/listing", app.DeveloperAlterListingHandler)
	dev.DELETE("/developer/listing", app.DeveloperDeleteListingHandler)

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
