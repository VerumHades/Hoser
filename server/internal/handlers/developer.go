package handlers

import (
	"net/http"
	"server/internal/database"

	"github.com/labstack/echo/v4"
)

// =================== TYPES ===================

type ApiDeveloperListing struct {
	Author      string                     `json:"author"`
	ID          string                     `json:"id"`
	Title       string                     `json:"title"`
	Description string                     `json:"description"`
	AccessMode  database.ListingAccessMode `json:"accessMode"`
}

type AddListingRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type ListingRequest struct {
	ID          string  `json:"id" validate:"required"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	AccessMode  *int    `json:"accessMode,omitempty"`
}

// =================== HELPERS ===================

func makeApiDeveloperListing(author string, l database.Listing) ApiDeveloperListing {
	return ApiDeveloperListing{
		Author:      author,
		ID:          l.GetUUID(),
		Title:       l.GetTitle(),
		Description: l.GetDescription(),
		AccessMode:  l.GetAccessMode(),
	}
}

// =================== HANDLERS ===================

// DeveloperListingsHandler returns all listings for the authenticated developer
func (app *App) DeveloperListingsHandler(c echo.Context) error {
	user, err := app.GetUserFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	dbListings, err := user.GetListings()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
	}

	listings := make([]ApiDeveloperListing, len(dbListings))
	username := user.GetUsername()
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

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	listing, err := user.CreateListing(req.Title, req.Description)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
	}

	return c.JSON(http.StatusOK, makeApiDeveloperListing(user.GetUsername(), listing))
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

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	listing, err := user.GetListing(req.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Listing not found")
	}

	if req.Title != nil {
		listing.SetTitle(*req.Title)
	}
	if req.Description != nil {
		listing.SetDescription(*req.Description)
	}
	if req.AccessMode != nil {
		mode, err := database.ListingAccessModeFromInt(*req.AccessMode)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid access mode")
		}
		listing.SetAccessMode(mode)
	}

	return c.JSON(http.StatusOK, makeApiDeveloperListing(user.GetUsername(), listing))
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

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := user.DeleteListing(req.ID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not delete listing")
	}

	return c.NoContent(http.StatusNoContent)
}
