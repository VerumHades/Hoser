package user

import (
	"api/internal/http/authentification"
	"api/pkg/util"
	"common/pkg/domain/user"
	"common/pkg/shared"
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type userLibraryQueryService interface {
	ExistsByUserAndListing(
		ctx context.Context,
		userID shared.UserID,
		listingID shared.ListingID,
	) (bool, error)

	FetchNextBatchOfUserSavedListings(
		ctx context.Context,
		userID shared.UserID,
		request shared.BatchRequest[user.UserSavedListingViewCursor],
	) (items []*user.UserSavedListingView, nextCursor user.UserSavedListingViewCursor, err error)
}

type userLibraryCommandService interface {
	SaveListingToLibrary(
		ctx context.Context,
		userID shared.UserID,
		listingID shared.ListingID,
	) (*user.SavedListing, error)

	RemoveListingFromLibrary(
		ctx context.Context,
		userID shared.UserID,
		saveID shared.SavedListingID,
	) error
}

type UserLibraryAPI struct {
	userLibraryQueryService   userLibraryQueryService
	userLibraryCommandService userLibraryCommandService
}

// --------------------
// Route registration
// --------------------

func (api *UserLibraryAPI) RegisterRoutes(group *echo.Group) {
	group.GET("/library", authentification.WithAuthenticatedUser(api.ListLibraryHandler))
	group.POST("/library", authentification.WithAuthenticatedUser(api.AddListingHandler))
	group.DELETE("/library", authentification.WithAuthenticatedUser(api.RemoveListingHandler))
	group.GET("/library/:listingId", authentification.WithAuthenticatedUser(api.HasListingHandler))
}

// --------------------
// API DTOs
// --------------------

type ApiListingBase struct {
	SavedListingID string  `json:"saved_listing_id"`
	ListingID      string  `json:"listing_id"`
	Title          *string `json:"title,omitempty"`
	Description    *string `json:"description,omitempty"`
	CreatedAt      string  `json:"created_at"`
}

// --------------------
// Conversion functions
// --------------------

func convertUserSavedListingViewToApi(domainListing *user.UserSavedListingView) ApiListingBase {
	return ApiListingBase{
		SavedListingID: string(domainListing.SavedListingID),
		ListingID:      string(domainListing.ListingID),
		Title:          domainListing.Title,
		Description:    domainListing.Description,
		CreatedAt:      domainListing.CreatedAt.Format(time.RFC3339),
	}
}

func convertBatchOfUserSavedListingViews(domainListings []*user.UserSavedListingView) []ApiListingBase {
	apiListings := make([]ApiListingBase, len(domainListings))
	for i, listing := range domainListings {
		apiListings[i] = convertUserSavedListingViewToApi(listing)
	}
	return apiListings
}

// --------------------
// Handler wrapper
// --------------------

// --------------------
// Handlers
// --------------------

func (api *UserLibraryAPI) ListLibraryHandler(userID shared.UserID, c echo.Context) error {
	ctx := c.Request().Context()
	batchSize := 50

	var cursor user.UserSavedListingViewCursor
	if encoded := c.QueryParam("cursor"); encoded != "" {
		if decoded, err := util.DecodeCursor[user.UserSavedListingViewCursor](encoded); err == nil {
			cursor = decoded
		}
	}

	batchRequest := shared.BatchRequest[user.UserSavedListingViewCursor]{
		Cursor:       cursor,
		MaxBatchSize: batchSize,
	}

	savedItems, nextCursor, err := api.userLibraryQueryService.FetchNextBatchOfUserSavedListings(ctx, userID, batchRequest)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch library")
	}

	apiListings := convertBatchOfUserSavedListingViews(savedItems)
	encodedCursor, _ := util.EncodeCursor(nextCursor)

	resp := util.PaginatedResponse[ApiListingBase, user.UserSavedListingViewCursor]{
		Items:  apiListings,
		Cursor: encodedCursor,
	}

	return c.JSON(http.StatusOK, resp)
}

func (api *UserLibraryAPI) AddListingHandler(userID shared.UserID, c echo.Context) error {
	ctx := c.Request().Context()

	var req struct {
		ListingID shared.ListingID `json:"id"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	_, err := api.userLibraryCommandService.SaveListingToLibrary(ctx, userID, req.ListingID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save listing to library")
	}

	return c.NoContent(http.StatusOK)
}

func (api *UserLibraryAPI) RemoveListingHandler(userID shared.UserID, c echo.Context) error {
	ctx := c.Request().Context()

	var req struct {
		ItemID shared.SavedListingID `json:"id"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	if err := api.userLibraryCommandService.RemoveListingFromLibrary(ctx, userID, req.ItemID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to remove listing from library")
	}

	return c.NoContent(http.StatusNoContent)
}

func (api *UserLibraryAPI) HasListingHandler(userID shared.UserID, c echo.Context) error {
	ctx := c.Request().Context()
	listingID := c.Param("listingId")
	if listingID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Listing ID is required")
	}

	exists, err := api.userLibraryQueryService.ExistsByUserAndListing(ctx, userID, shared.ListingID(listingID))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch user library")
	}

	return c.JSON(http.StatusOK, map[string]bool{"hasListing": exists})
}
