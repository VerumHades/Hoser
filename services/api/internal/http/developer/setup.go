package developer

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

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
