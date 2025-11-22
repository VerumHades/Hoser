package handlers

import (
	"fmt"
	"net/http"
	"server/internal/database"

	"github.com/labstack/echo/v4"
)

// =================== TYPES ===================

type ApiDeveloperListing struct {
	ID          string          `json:"id"`
	Title       *string         `json:"title,omitempty"`
	Description *string         `json:"description,omitempty"`
	AccessMode  *int            `json:"accessMode,omitempty"`
	Prices      *PricesRequest  `json:"prices,omitempty"`
	Hardware    *HardwareUpdate `json:"hardware,omitempty"`
	Author      string          `json:"author"` // kept for reference
}
type AddListingRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ListingRequest struct {
	ID          string          `json:"id"`
	Title       *string         `json:"title,omitempty"`
	Description *string         `json:"description,omitempty"`
	AccessMode  *int            `json:"accessMode,omitempty"`
	Prices      *PricesRequest  `json:"prices,omitempty"`
	Hardware    *HardwareUpdate `json:"hardware,omitempty"`
}

type PricesRequest struct {
	SinglePurchase      *CurrencyRequest `json:"singlePurchase,omitempty"`
	MonthlySubscription *CurrencyRequest `json:"monthlySubscription,omitempty"`
	MonthlyHardware     *CurrencyRequest `json:"monthlyHardware,omitempty"`
}

type CurrencyRequest struct {
	Value float32 `json:"value"`
	Name  string  `json:"name"`
	Short string  `json:"short"`
}

type HardwareUpdate struct {
	CPU  *int   `json:"cpu,omitempty"`
	RAM  *int64 `json:"ramBytes,omitempty"`
	Disk *int64 `json:"diskBytes,omitempty"`
}

// =================== HELPERS ===================

func makeApiDeveloperListing(author string, l database.Listing) ApiDeveloperListing {
	// Map hardware requirements
	hwSpec := l.HardwareRequirements()
	cpu := hwSpec.CPUCount()
	ram := hwSpec.RAMBytes()
	disk := hwSpec.DiskBytes()
	hardware := &HardwareUpdate{
		CPU:  &cpu,
		RAM:  &ram,
		Disk: &disk,
	}

	// Map prices
	prices := &PricesRequest{
		SinglePurchase:      mapCurrency(l.SinglePurchasePrice()),
		MonthlySubscription: mapCurrency(l.MonthlySubscriptionPrice()),
		MonthlyHardware:     mapCurrency(l.HardwareRequirements().MonthlyPrice()),
	}

	// Map access mode
	access := int(l.AccessMode())

	// Title & description
	title := l.Title()
	description := l.Description()

	listing := ApiDeveloperListing{
		Author:      author,
		ID:          l.UUID(),
		Title:       &title,
		Description: &description,
		AccessMode:  &access,
		Prices:      prices,
		Hardware:    hardware,
	}

	fmt.Println(listing)

	return listing
}

// Map Currency interface to *CurrencyRequest
func mapCurrency(c database.Currency) *CurrencyRequest {
	return &CurrencyRequest{
		Value: c.AsNumber(),
		Name:  c.Name(),
		Short: c.Short(),
	}
}

func parseAccessMode(i int) (database.ListingAccessMode, error) {
	switch database.ListingAccessMode(i) {
	case database.Private, database.Public:
		return database.ListingAccessMode(i), nil
	default:
		return database.Private, echo.NewHTTPError(http.StatusBadRequest, "Invalid access mode")
	}
}

// =================== HANDLERS ===================

// DeveloperListingsHandler returns all listings for the authenticated developer
func (app *App) DeveloperListingsHandler(c echo.Context) error {
	user, err := app.GetUserFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	dbListings, err := user.Listings()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
	}

	listings := make([]ApiDeveloperListing, len(dbListings))
	username := user.Username()

	for i, l := range dbListings {
		listings[i] = makeApiDeveloperListing(username, l)
	}

	return c.JSON(http.StatusOK, listings)
}

// DeveloperAddListingHandler adds a new listing for the developer
func (app *App) DeveloperAddListingHandler(c echo.Context) error {
	user, err := app.GetUserFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	var req AddListingRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	listing, err := user.CreateListing(req.Title, req.Description)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
	}

	return c.JSON(http.StatusOK, makeApiDeveloperListing(user.Username(), listing))
}

// DeveloperAlterListingHandler modifies an existing listing
func (app *App) DeveloperAlterListingHandler(c echo.Context) error {
	user, err := app.GetUserFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	var req ListingRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	listing, err := user.Listing(req.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Listing not found")
	}
	// Update title/description/access mode
	if req.Title != nil {
		if err := listing.SetTitle(*req.Title); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}
	if req.Description != nil {
		if err := listing.SetDescription(*req.Description); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}
	if req.AccessMode != nil {
		mode, err := parseAccessMode(*req.AccessMode)
		if err != nil {
			return err
		}
		if err := listing.SetAccessMode(mode); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}

	if req.Prices != nil {
		if req.Prices.SinglePurchase != nil {
			cp := listing.SinglePurchasePrice()
			cp.SetValue(req.Prices.SinglePurchase.Value)
			cp.SetName(req.Prices.SinglePurchase.Name)
			cp.SetShort(req.Prices.SinglePurchase.Short)
		}
		if req.Prices.MonthlySubscription != nil {
			cp := listing.MonthlySubscriptionPrice()
			cp.SetValue(req.Prices.MonthlySubscription.Value)
			cp.SetName(req.Prices.MonthlySubscription.Name)
			cp.SetShort(req.Prices.MonthlySubscription.Short)
		}
		if req.Prices.MonthlyHardware != nil {
			cp := listing.HardwareRequirements().MonthlyPrice()
			cp.SetValue(req.Prices.MonthlyHardware.Value)
			cp.SetName(req.Prices.MonthlyHardware.Name)
			cp.SetShort(req.Prices.MonthlyHardware.Short)
		}
	}

	// Update hardware requirements
	if req.Hardware != nil {
		hw := listing.HardwareRequirements()
		if req.Hardware.CPU != nil {
			_ = hw.SetCPUCount(*req.Hardware.CPU)
		}
		if req.Hardware.RAM != nil {
			_ = hw.SetRAMBytes(*req.Hardware.RAM)
		}
		if req.Hardware.Disk != nil {
			_ = hw.SetDiskBytes(*req.Hardware.Disk)
		}
	}

	return c.JSON(http.StatusOK, makeApiDeveloperListing(user.Username(), listing))
}

// DeveloperDeleteListingHandler deletes a listing
func (app *App) DeveloperDeleteListingHandler(c echo.Context) error {
	user, err := app.GetUserFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	var req ListingRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	if err := user.DeleteListing(req.ID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not delete listing")
	}

	return c.NoContent(http.StatusNoContent)
}
