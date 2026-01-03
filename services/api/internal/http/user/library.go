package userapi

import (
	"api/internal/http/authentification"
	"api/pkg/util"
	"common/pkg/domain/entities/user"
	"common/pkg/shared"
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type APIUserLibraryQueryService interface {
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

type APIUserLibraryCommandService interface {
	SaveListingToLibrary(
		ctx context.Context,
		userID shared.UserID,
		listingID shared.ListingID,
	) (*user.SavedListing, error)

	RemoveListingFromLibrary(
		ctx context.Context,
		userID shared.UserID,
		listingID shared.ListingID,
	) error
}

type UserLibraryAPI struct {
	userLibraryQueryService   APIUserLibraryQueryService
	userLibraryCommandService APIUserLibraryCommandService
}

// NewUserLibraryAPI constructs a UserLibraryAPI.
func NewUserLibraryAPI(
	userLibraryQueryService APIUserLibraryQueryService,
	userLibraryCommandService APIUserLibraryCommandService,
) *UserLibraryAPI {
	return &UserLibraryAPI{
		userLibraryQueryService:   userLibraryQueryService,
		userLibraryCommandService: userLibraryCommandService,
	}
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
	ListingID      string  `json:"id"`
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

// --------------------
// Handler wrapper
// --------------------

// --------------------
// Handlers
// --------------------

func (api *UserLibraryAPI) ListLibraryHandler(userID shared.UserID, c echo.Context) error {
	return util.HandleBatchRequest(
		c,
		func(
			ctx context.Context,
			request shared.BatchRequest[user.UserSavedListingViewCursor],
		) (items []*user.UserSavedListingView, nextCursor user.UserSavedListingViewCursor, err error) {
			return api.userLibraryQueryService.FetchNextBatchOfUserSavedListings(ctx, userID, request)
		},
		func(elements []*user.UserSavedListingView) (views []ApiListingBase) {
			return util.MapList(elements, convertUserSavedListingViewToApi)
		},
	)
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
		ItemID shared.ListingID `json:"id"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	if err := api.userLibraryCommandService.RemoveListingFromLibrary(ctx, userID, req.ItemID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to remove listing from library:"+err.Error())
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
