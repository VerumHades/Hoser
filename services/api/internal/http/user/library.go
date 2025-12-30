package user

import (
	"api/pkg/util"
	"common/pkg/domain/user"
	"common/pkg/shared"
	"net/http"

	"github.com/labstack/echo/v4"
)

// =================== TYPES ===================
type ApiListingBase struct {
	ID          string  `json:"id"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
}

// UserLibraryHandler returns all saved listings for the authenticated user.
func (api *UserAPI) UserLibraryHandler(c echo.Context) error {
	userID, err := GetAuthenticatedUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	ctx := c.Request().Context()
	batchSize := 50

	// Decode cursor from query param if present
	var cursor user.UserSavedListingViewCursor
	if encoded := c.QueryParam("cursor"); encoded != "" {
		if decodedCursor, err := util.DecodeCursor[user.UserSavedListingViewCursor](encoded); err == nil {
			cursor = decodedCursor
		}
	}

	// Create BatchRequest with typed cursor
	batchRequest := shared.BatchRequest[user.UserSavedListingViewCursor]{
		Cursor:       cursor,
		MaxBatchSize: batchSize,
	}

	// Fetch batch of saved listings
	savedItems, nextCursor, err := api.userLibraryQueryService.FetchNextBatchOfUserSavedListings(ctx, userID, batchRequest)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch library")
	}

	// Convert to API response objects
	apiListings := make([]ApiListingBase, 0, len(savedItems))
	for _, savedItem := range savedItems {
		user, err := app.ListingService.GetPublicListing(savedItem.ListingID)
		if err != nil {
			continue // skip missing or private listings
		}
		apiListings = append(apiListings, app.MakeApiDeveloperListing(user).ApiListingBase)
	}

	// Encode cursor for next request
	encodedCursor, _ := util.EncodeCursor(nextCursor)

	resp := PaginatedResponse[ApiListingBase, user.SavedListingCursor]{
		Items:  apiListings,
		Cursor: encodedCursor,
	}

	return c.JSON(http.StatusOK, resp)
}

// UserAddListingToLibraryHandler saves a user to the authenticated user's library
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
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save user to library")
	}

	return c.NoContent(http.StatusOK)
}

// UserRemoveListingFromLibraryHandler removes a saved user from the authenticated user's library
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
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to remove user from library")
	}

	return c.NoContent(http.StatusNoContent)
}

func (app *App) HasListingInLibraryHandler(c echo.Context) error {
	// Extract user ID from JWT
	userID, err := app.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	// Get user ID from URL path
	listingID := c.Param("listingId")
	if listingID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Listing ID is required")
	}

	// List the user's library
	libraryItems, err := app.UserAppService.ListUserLibrary(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch user library")
	}

	// Check if the user exists in the library
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
