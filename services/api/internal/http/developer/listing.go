package handlers

import (
	"common/pkg/domain/listing"
	"common/pkg/shared"
	"context"

	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type listingService interface {
	CreateListing(ctx context.Context, listing *listing.Listing) error
	UpdateListing(ctx context.Context, listingID shared.ListingID, updateFunction func(listing *listing.Listing) error) error
	DeleteListing(ctx context.Context, listingID shared.ListingID) error
}

type DeveloperListingAPI struct {
	listingService listingService
}

// =================== TYPES ===================
type ApiListingBase struct {
	ID          string                        `json:"id"`
	Title       *string                       `json:"title,omitempty"`
	Description *string                       `json:"description,omitempty"`
	Price       int64                         `json:"price"`
	Hardware    *shared.HardwareSpecification `json:"hardware,omitempty"`
	Author      string                        `json:"author"`
}

type ApiDeveloperListing struct {
	ApiListingBase
	AccessMode *int `json:"accessMode,omitempty"`
}

type ApiPublicListing struct {
	ApiListingBase
}

type UpdateListingRequest struct {
	ID          string                        `json:"id"`
	Title       *string                       `json:"title,omitempty"`
	Description *string                       `json:"description,omitempty"`
	Hardware    *shared.HardwareSpecification `json:"hardware,omitempty"`
	Price       *int64                        `json:"price,omitempty"`
	AccessMode  *listing.ListingAccessMode    `json:"accessMode,omitempty"` // integer type
}

type AddListingRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// =================== HELPERS ===================
func (app *App) MakeApiDeveloperListing(l *listing.Listing) ApiDeveloperListing {
	accessMode := int(l.AccessMode) // convert ListingAccessMode to int

	return ApiDeveloperListing{
		ApiListingBase: ApiListingBase{
			ID:          l.ID,
			Author:      l.AuthorID,
			Title:       &l.Title,
			Description: &l.Description,
			Price:       l.Price,
			Hardware:    l.HardwareSpecification,
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

	listing, err := app.ListingService.UpdateListing(userID, req.ID, listing.ListingUpdate{
		Title:                 req.Title,
		Description:           req.Description,
		HardwareSpecification: req.Hardware,
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

func (app *App) HasListingInLibraryHandler(c echo.Context) error {
	// Extract user ID from JWT
	userID, err := app.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	// Get listing ID from URL path
	listingID := c.Param("listingId")
	if listingID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Listing ID is required")
	}

	// List the user's library
	libraryItems, err := app.UserAppService.ListUserLibrary(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch user library")
	}

	// Check if the listing exists in the library
	hasListing := false
	for _, item := range libraryItems {
		if item.ListingID == listingID {
			hasListing = true
			break
		}
	}

	// Return result
	return c.JSON(http.StatusOK, map[string]bool{
		"hasListing": hasListing,
	})
}

// =================== GITHUB SETUP HANDLERS ===================

type GithubSetupResponse struct {
	RepoURL     string `json:"repoUrl"`
	AccessToken string `json:"accessToken,omitempty"` // optional, maybe masked
}

// DeveloperGetSetupHandler returns the GitHub setup for a listing owned by the user
func (app *App) DeveloperGetSetupHandler(c echo.Context) error {
	userID, err := app.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	listingID := c.Param("listingId")
	if listingID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Listing ID is required")
	}

	setup, err := app.ListingService.GetSetupForListingPrivate(userID, listingID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "GitHub setup not found")
	}

	// Only expose safe fields
	response := GithubSetupResponse{
		RepoURL:     setup.RepoURL,
		AccessToken: "", // never send the raw token
	}

	return c.JSON(http.StatusOK, response)
}

// AttachOrUpdateSetupHandler attaches or updates a GitHub setup for a listing
func (app *App) DeveloperAttachOrUpdateSetupHandler(c echo.Context) error {
	userID, err := app.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	listingID := c.Param("listingId")
	if listingID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Listing ID is required")
	}

	var req struct {
		RepoURL     string `json:"repoUrl"`
		AccessToken string `json:"accessToken"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	setup, err := app.ListingService.AttachOrUpdateSetup(userID, listingID, req.RepoURL, req.AccessToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to attach/update setup %s", err))
	}

	return c.JSON(http.StatusOK, setup)
}

// RemoveSetupHandler removes the GitHub setup for a listing
func (app *App) DeveloperRemoveSetupHandler(c echo.Context) error {
	userID, err := app.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	listingID := c.Param("listingId")
	if listingID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Listing ID is required")
	}

	if err := app.ListingService.RemoveSetup(userID, listingID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to remove setup")
	}

	return c.NoContent(http.StatusNoContent)
}
