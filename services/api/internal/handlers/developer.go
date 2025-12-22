package handlers

import (
	"common/pkg/hardware"
	"common/pkg/listing"
	"common/pkg/money"
	"net/http"

	"github.com/labstack/echo/v4"
)

// =================== TYPES ===================
type ApiListingBase struct {
	ID          string                    `json:"id"`
	Title       *string                   `json:"title,omitempty"`
	Description *string                   `json:"description,omitempty"`
	Price       CurrencyRequest           `json:"price"`
	Hardware    *HardwareSpecificationDTO `json:"hardware,omitempty"`
	Author      string                    `json:"author"`
}

type ApiDeveloperListing struct {
	ApiListingBase
	AccessMode *int `json:"accessMode,omitempty"`
}

type ApiPublicListing struct {
	ApiListingBase
}

type UpdateListingRequest struct {
	ID          string                     `json:"id"`
	Title       *string                    `json:"title,omitempty"`
	Description *string                    `json:"description,omitempty"`
	Hardware    *HardwareSpecificationDTO  `json:"hardware,omitempty"`
	Price       *money.Money               `json:"price,omitempty"`
	AccessMode  *listing.ListingAccessMode `json:"accessMode,omitempty"` // integer type
}
type ListingRequest struct {
	ID          string                    `json:"id"`
	Title       *string                   `json:"title,omitempty"`
	Description *string                   `json:"description,omitempty"`
	AccessMode  *int                      `json:"accessMode,omitempty"`
	Price       *CurrencyRequest          `json:"price,omitempty"`
	Hardware    *HardwareSpecificationDTO `json:"hardware,omitempty"`
}

type AddListingRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// =================== HELPERS ===================
func (app *App) MakeApiDeveloperListing(l *listing.Listing) ApiDeveloperListing {
	accessMode := int(l.AccessMode) // convert ListingAccessMode to int

	price := CurrencyRequest{
		Value: float32(l.Price.Amount),
		Short: l.Price.CurrencyCode,
	}

	return ApiDeveloperListing{
		ApiListingBase: ApiListingBase{
			ID:          l.ID,
			Author:      l.AuthorID,
			Title:       &l.Title,
			Description: &l.Description,
			Price:       price,
			Hardware:    app.HardwareSpecificationToDTO(l.HardwareSpecification),
		},
		AccessMode: &accessMode,
	}
}

func DeveloperToPublicListing(dev ApiDeveloperListing) ApiPublicListing {
	return ApiPublicListing{
		ApiListingBase: dev.ApiListingBase,
	}
}

// DeveloperListingsHandler returns all listings for the authenticated developer
func (app *App) DeveloperListingsHandler(c echo.Context) error {
	userID, err := app.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	dbListings, err := app.ListingService.ListByAuthor(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
	}

	listings := make([]ApiDeveloperListing, len(dbListings))

	for i, l := range dbListings {
		listings[i] = app.MakeApiDeveloperListing(l)
	}

	return c.JSON(http.StatusOK, listings)
}

// DeveloperGetListingHandler returns a single listing for the authenticated developer by ID
func (app *App) DeveloperGetListingHandler(c echo.Context) error {
	userID, err := app.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	listingID := c.Param("id")
	if listingID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Listing ID is required")
	}

	listing, err := app.ListingService.GetOwnedListing(userID, listingID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Listing not found")
	}

	apiListing := app.MakeApiDeveloperListing(listing)
	return c.JSON(http.StatusOK, apiListing)
}

// DeveloperAddListingHandler adds a new listing for the developer
func (app *App) DeveloperAddListingHandler(c echo.Context) error {
	userID, err := app.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	var req AddListingRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	listing, err := app.UserAppService.CreateListing(userID, req.Title, req.Description, listing.Private)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
	}

	return c.JSON(http.StatusOK, app.MakeApiDeveloperListing(listing))
}

// DeveloperUpdateListingHandler updates an existing listing for the developer
func (app *App) DeveloperUpdateListingHandler(c echo.Context) error {
	userID, err := app.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	var req UpdateListingRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	// Translate HardwareUpdate to HardwareSpecification
	var hardwareSpec hardware.HardwareSpecification
	if req.Hardware != nil {
		if req.Hardware.CPU != nil {
			hardwareSpec.CPUCount = *req.Hardware.CPU
		}
		if req.Hardware.RAM != nil {
			hardwareSpec.RAMBytes = *req.Hardware.RAM
		}
		if req.Hardware.Disk != nil {
			hardwareSpec.DiskBytes = *req.Hardware.Disk
		}
	}

	listing, err := app.ListingService.UpdateListing(userID, req.ID, listing.ListingUpdate{
		Title:                 req.Title,
		Description:           req.Description,
		HardwareSpecification: &hardwareSpec,
		Price:                 req.Price,
		AccessMode:            req.AccessMode,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update listing")
	}

	return c.JSON(http.StatusOK, app.MakeApiDeveloperListing(listing))
}

// DeveloperDeleteListingHandler deletes a listing
func (app *App) DeveloperDeleteListingHandler(c echo.Context) error {
	userID, err := app.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	var req ListingRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	if err := app.ListingService.DeleteListing(userID, req.ID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not delete listing")
	}

	return c.NoContent(http.StatusNoContent)
}

// UserLibraryHandler returns all saved listings for the authenticated user.
func (app *App) UserLibraryHandler(c echo.Context) error {
	userID, err := app.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	savedItems, err := app.UserAppService.ListUserLibrary(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch library")
	}

	apiListings := make([]ApiListingBase, 0, len(savedItems))
	for _, savedItem := range savedItems {
		listing, err := app.ListingService.GetPublicListing(savedItem.ListingID)
		if err != nil {
			continue // skip listings that no longer exist or are not public
		}
		apiListings = append(apiListings, app.MakeApiDeveloperListing(listing).ApiListingBase)
	}

	return c.JSON(http.StatusOK, apiListings)
}

// UserAddListingToLibraryHandler saves a listing to the authenticated user's library
func (app *App) UserAddListingToLibraryHandler(c echo.Context) error {
	userID, err := app.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	var req struct {
		ListingID string `json:"id"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	_, err = app.UserAppService.SaveListingToLibrary(userID, req.ListingID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save listing to library")
	}

	return c.NoContent(http.StatusOK)
}

// UserRemoveListingFromLibraryHandler removes a saved listing from the authenticated user's library
func (app *App) UserRemoveListingFromLibraryHandler(c echo.Context) error {
	userID, err := app.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	var req struct {
		ItemID string `json:"id"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	if err := app.UserAppService.RemoveListingFromLibrary(userID, req.ItemID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to remove listing from library")
	}

	return c.NoContent(http.StatusNoContent)
}
