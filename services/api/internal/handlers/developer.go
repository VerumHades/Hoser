package handlers

import (
	"api/internal/database"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// =================== TYPES ===================
type ApiListingBase struct {
	ID          string            `json:"id"`
	Title       *string           `json:"title,omitempty"`
	Description *string           `json:"description,omitempty"`
	Prices      []ApiPricingEntry `json:"prices"`
	Hardware    *HardwareUpdate   `json:"hardware,omitempty"`
	Author      string            `json:"author"`
}

type ApiDeveloperListing struct {
	ApiListingBase
	AccessMode *int `json:"accessMode,omitempty"`
}

type ApiPublicListing struct {
	ApiListingBase
}

type ApiPricingEntry struct {
	ID       string          `json:"id"`   // "0", "1", etc
	Type     int             `json:"type"` // 0 = OneTime, 1 = Monthly
	Currency CurrencyRequest `json:"currency"`
}

type ListingRequest struct {
	ID          string            `json:"id"`
	Title       *string           `json:"title,omitempty"`
	Description *string           `json:"description,omitempty"`
	AccessMode  *int              `json:"accessMode,omitempty"`
	Prices      []ApiPricingEntry `json:"prices,omitempty"`
	Hardware    *HardwareUpdate   `json:"hardware,omitempty"`
}

type AddListingRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
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

func (app *App) MakeApiDeveloperListing(l database.Listing) ApiDeveloperListing {
	// Hardware
	hwSpec := l.HardwareRequirements()
	cpu := hwSpec.CPUCount()
	ram := hwSpec.RAMBytes()
	disk := hwSpec.DiskBytes()
	hardware := &HardwareUpdate{CPU: &cpu, RAM: &ram, Disk: &disk}

	priceEntries := app.ConvertPricingListToAPI(l.Pricing())
	fmt.Println(priceEntries)
	// Access mode
	access := int(l.AccessMode())
	title := l.Title()
	description := l.Description()

	return ApiDeveloperListing{
		ApiListingBase: ApiListingBase{
			ID:          l.UUID(),
			Author:      l.Author().Username(),
			Title:       &title,
			Description: &description,
			Prices:      priceEntries,
			Hardware:    hardware,
		},
		AccessMode: &access,
	}
}
func DeveloperToPublicListing(dev ApiDeveloperListing) ApiPublicListing {
	return ApiPublicListing{
		ApiListingBase: dev.ApiListingBase,
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

func (app *App) PublicGetListingHandler(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "listing id required")
	}

	listing, err := app.DatabaseInteractor.GetPublicListing(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Listing not found")
	}

	return c.JSON(http.StatusOK, app.MakeApiDeveloperListing(listing).ApiListingBase)
}

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

	for i, l := range dbListings {
		listings[i] = app.MakeApiDeveloperListing(l)
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

	return c.JSON(http.StatusOK, app.MakeApiDeveloperListing(listing))
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

	return c.JSON(http.StatusOK, app.MakeApiDeveloperListing(listing))
}

type AddPricingRequest struct {
	ListingID string          `json:"listingId"`
	Type      int             `json:"type"` // 0 = OneTime, 1 = Monthly
	Currency  CurrencyRequest `json:"currency"`
}

type UpdatePricingRequest struct {
	ListingID string          `json:"listingId"`
	PricingID string          `json:"pricingId"`
	Type      int             `json:"type"`
	Currency  CurrencyRequest `json:"currency"`
}

type DeletePricingRequest struct {
	ListingID string `json:"listingId"`
	PricingID string `json:"pricingId"`
}

func (app *App) DeveloperAddListingPricingHandler(c echo.Context) error {
	user, err := app.GetUserFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	var req AddPricingRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}
	if req.ListingID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "listingId required")
	}

	listing, err := user.Listing(req.ListingID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Listing not found")
	}

	pl := listing.Pricing()
	newPricing, err := pl.AddPricing(
		database.PricingType(req.Type),
		app.DatabaseInteractor.NewCurrency(req.Currency.Name, req.Currency.Short, req.Currency.Value))

	return c.JSON(http.StatusCreated, ApiPricingEntry{
		ID:   newPricing.UUID(),
		Type: int(newPricing.Type()),
		Currency: CurrencyRequest{
			Value: newPricing.Amount().AsNumber(),
			Name:  newPricing.Amount().Name(),
			Short: newPricing.Amount().Short(),
		},
	})
}

func (app *App) DeveloperUpdateListingPricingHandler(c echo.Context) error {
	user, err := app.GetUserFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	var req UpdatePricingRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	if req.ListingID == "" || req.PricingID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "listingId and pricingId required")
	}

	listing, err := user.Listing(req.ListingID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Listing not found")
	}

	pl := listing.Pricing()
	p := pl.GetPricing(req.PricingID)
	if p == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Pricing not found")
	}

	// update the pricing
	p.SetType(database.PricingType(req.Type))
	p.Amount().SetValue(req.Currency.Value)
	p.Amount().SetName(req.Currency.Name)
	p.Amount().SetShort(req.Currency.Short)

	// return updated entry
	return c.JSON(http.StatusOK, ApiPricingEntry{
		ID:   req.PricingID,
		Type: int(p.Type()),
		Currency: CurrencyRequest{
			Value: p.Amount().AsNumber(),
			Name:  p.Amount().Name(),
			Short: p.Amount().Short(),
		},
	})
}

func (app *App) DeveloperDeleteListingPriceHandler(c echo.Context) error {
	user, err := app.GetUserFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	var req DeletePricingRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	listing, err := user.Listing(req.ListingID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Listing not found")
	}

	err = listing.Pricing().RemovePricing(req.PricingID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Pricing not found")
	}

	return c.NoContent(http.StatusNoContent)
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
