package userapi

import (
	"api/pkg/util"
	"common/pkg/domain/listing"
	"common/pkg/shared"
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type APIPublicUserListingQueryService interface {
	GetListing(
		ctx context.Context,
		listingID shared.ListingID,
	) (*listing.Listing, error)

	FetchNextListingsByAuthor(
		ctx context.Context,
		authorID shared.UserID,
		request shared.BatchRequest[listing.ListingCursor],
	) ([]*listing.Listing, listing.ListingCursor, error)

	SearchNextListingsBatch(
		ctx context.Context,
		query string,
		request shared.BatchRequest[listing.ListingSearchCursor],
	) ([]*listing.Listing, listing.ListingSearchCursor, error)
}

type PublicUserAPI struct {
	publicUserListingQueryService APIPublicUserListingQueryService
}

// NewPublicUserAPI constructs a PublicUserAPI
func NewPublicUserAPI(
	publicUserListingQueryService APIPublicUserListingQueryService,
) *PublicUserAPI {
	return &PublicUserAPI{
		publicUserListingQueryService: publicUserListingQueryService,
	}
}

// --------------------
// Route registration
// --------------------
func (api *PublicUserAPI) RegisterRoutes(group *echo.Group) {
	group.GET("/listings/:listingId", api.GetListingHandler)
	//group.GET("/authors/:authorId/listings", api.ListAuthorListingsHandler)
	group.GET("/search/listings", api.SearchListingsHandler)
}

// --------------------
// API DTOs
// --------------------
type ApiListing struct {
	ID          string `json:"id"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	AuthorID    string `json:"author_id"`
	CreatedAt   string `json:"created_at"`
}

// --------------------
// Conversion functions
// --------------------
func convertListingToApi(domainListing *listing.Listing) ApiListing {
	return ApiListing{
		ID:          string(domainListing.ID()),
		Title:       domainListing.Title(),
		Description: domainListing.Description(),
		CreatedAt:   domainListing.CreatedAt().Format(time.RFC3339),
	}
}

// --------------------
// Handlers
// --------------------
func (api *PublicUserAPI) GetListingHandler(c echo.Context) error {
	ctx := c.Request().Context()
	listingID := c.Param("listingId")
	if listingID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Listing ID is required")
	}

	listing, err := api.publicUserListingQueryService.GetListing(ctx, shared.ListingID(listingID))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch listing")
	}
	if listing == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Listing not found")
	}

	return c.JSON(http.StatusOK, convertListingToApi(listing))
}

func (api *PublicUserAPI) ListAuthorListingsHandler(c echo.Context) error {
	//ctx := c.Request().Context()
	authorID := c.Param("authorId")
	if authorID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Author ID is required")
	}

	return util.HandleBatchRequest(
		c,
		func(ctx context.Context, request shared.BatchRequest[listing.ListingCursor]) ([]*listing.Listing, listing.ListingCursor, error) {
			return api.publicUserListingQueryService.FetchNextListingsByAuthor(ctx, shared.UserID(authorID), request)
		},
		func(elements []*listing.Listing) (views []ApiListing) {
			return util.MapList(elements, convertListingToApi)
		},
	)
}

func (api *PublicUserAPI) SearchListingsHandler(c echo.Context) error {
	//ctx := c.Request().Context()
	query := c.QueryParam("q")
	//if query == "" {
	//	return echo.NewHTTPError(http.StatusBadRequest, "Search query is required")
	//}

	return util.HandleBatchRequest(
		c,
		func(ctx context.Context, request shared.BatchRequest[listing.ListingSearchCursor]) ([]*listing.Listing, listing.ListingSearchCursor, error) {
			return api.publicUserListingQueryService.SearchNextListingsBatch(ctx, query, request)
		},
		func(elements []*listing.Listing) (views []ApiListing) {
			return util.MapList(elements, convertListingToApi)
		},
	)
}
