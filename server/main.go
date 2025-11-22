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
)

type PublicListing struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Hardware    *HardwareDTO   `json:"hardware,omitempty"`
	Prices      *ListingPrices `json:"prices,omitempty"`
}

type HardwareDTO struct {
	CPU  int   `json:"cpu"`
	RAM  int64 `json:"ram"`
	Disk int64 `json:"disk"`
}

type ListingPrices struct {
	SinglePurchase      *CurrencyDTO `json:"singlePurchase,omitempty"`
	MonthlySubscription *CurrencyDTO `json:"monthlySubscription,omitempty"`
	MonthlyHardware     *CurrencyDTO `json:"monthlyHardware,omitempty"`
}

type CurrencyDTO struct {
	Name  string  `json:"name"`
	Short string  `json:"short"`
	Value float32 `json:"value"`
}

func ConvertListingToPublic(l database.Listing) *PublicListing {
	return &PublicListing{
		ID:          l.UUID(),
		Title:       l.Title(),
		Description: l.Description(),
		Hardware: &HardwareDTO{
			CPU:  l.HardwareRequirements().CPUCount(),
			RAM:  l.HardwareRequirements().RAMBytes(),
			Disk: l.HardwareRequirements().DiskBytes(),
		},
		Prices: &ListingPrices{
			SinglePurchase:      convertCurrency(l.SinglePurchasePrice()),
			MonthlySubscription: convertCurrency(l.MonthlySubscriptionPrice()),
			MonthlyHardware:     convertCurrency(l.HardwareRequirements().MonthlyPrice()),
		},
	}
}

func convertCurrency(c database.Currency) *CurrencyDTO {
	if c == nil {
		return nil
	}
	return &CurrencyDTO{
		Name:  c.Name(),
		Short: c.Short(),
		Value: c.AsNumber(),
	}
}

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
	public.GET("/rentals/public", func(c echo.Context) error {
		options := database.ListingQueryOptions{}
		if q := c.QueryParam("q"); q != "" {
			options.Text = q
		}

		listings, err := app.DatabaseInteractor.QueryPublicListings(&options)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		// Convert to API response DTOs
		response := make([]*PublicListing, 0, len(listings))
		for _, l := range listings {
			response = append(response, ConvertListingToPublic(l))
		}

		return c.JSON(http.StatusOK, response)
	})
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
	dev := e.Group("/developer")
	dev.Use(app.DeveloperOnlyMiddleware)

	dev.GET("/listings", app.DeveloperListingsHandler)
	dev.POST("/listing", app.DeveloperAddListingHandler)
	dev.PUT("/listing", app.DeveloperAlterListingHandler)
	dev.DELETE("/listing", app.DeveloperDeleteListingHandler)

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
