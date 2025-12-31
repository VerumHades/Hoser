package developerapi

import (
	"api/internal/http/authentification"
	"api/pkg/util"
	"common/pkg/domain/listing"
	"common/pkg/shared"
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type listingService interface {
	CreateListing(ctx context.Context, l *listing.Listing) error
	UpdateListing(
		ctx context.Context,
		listingID shared.ListingID,
		updateFunc func(l *listing.Listing) error,
	) (l *listing.Listing, err error)
	DeleteListing(ctx context.Context, listingID shared.ListingID) error

	FetchNextBatchByAuthor(
		ctx context.Context,
		authorID shared.UserID,
		request shared.BatchRequest[listing.ListingCursor],
	) (listings []*listing.Listing, nextCursor listing.ListingCursor, err error)

	GetOwnedListing(ctx context.Context, listingID shared.ListingID, userID shared.UserID) (*listing.Listing, error)
}

type userService interface {
	IsUserDeveloper(ctx context.Context, userID shared.UserID) (bool, error)
}

// DeveloperListingAPI is the API layer for developer-specific listing endpoints
type DeveloperListingAPI struct {
	listingService listingService
	userService    userService
}

// NewDeveloperListingAPI constructs a DeveloperListingAPI with required dependencies.
func NewDeveloperListingAPI(
	listingService listingService,
	userService userService,
) *DeveloperListingAPI {
	return &DeveloperListingAPI{
		listingService: listingService,
		userService:    userService,
	}
}

// --------------------
// Route registration
// --------------------

func (api *DeveloperListingAPI) RegisterRoutes(group *echo.Group) {
	group.GET("/listings", authentification.WithAuthenticatedUser(api.ListListingsHandler))
	group.GET("/listings/:id", authentification.WithAuthenticatedUser(api.GetListingHandler))
	group.POST("/listings", authentification.WithAuthenticatedUser(api.AddListingHandler))
	group.PUT("/listings", authentification.WithAuthenticatedUser(api.UpdateListingHandler))
	group.DELETE("/listings/:id", authentification.WithAuthenticatedUser(api.DeleteListingHandler))
}

// --------------------
// API DTOs
// --------------------

type ApiListingBase struct {
	ID          shared.ListingID              `json:"id"`
	Title       string                        `json:"title,omitempty"`
	Description string                        `json:"description,omitempty"`
	Price       int64                         `json:"price"`
	Hardware    *shared.HardwareSpecification `json:"hardware,omitempty"`
}

type ApiDeveloperListing struct {
	ApiListingBase
	AccessMode *int `json:"accessMode,omitempty"`
}

type ApiPublicListing struct {
	ApiListingBase
}

type AddListingRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateListingRequest struct {
	ID          shared.ListingID              `json:"id"`
	Title       *string                       `json:"title,omitempty"`
	Description *string                       `json:"description,omitempty"`
	Hardware    *shared.HardwareSpecification `json:"hardware,omitempty"`
	Price       *int64                        `json:"price,omitempty"`
	AccessMode  *listing.ListingAccessMode    `json:"accessMode,omitempty"`
}

// --------------------
// Conversion helpers
// --------------------

func MakeApiDeveloperListing(l *listing.Listing) ApiDeveloperListing {
	accessMode := int(l.AccessMode())
	return ApiDeveloperListing{
		ApiListingBase: ApiListingBase{
			ID:          l.ID(),
			Title:       l.Title(),
			Description: l.Description(),
			Price:       l.PriceInMinorUnits(),
			Hardware:    l.HardwareSpecification(),
		},
		AccessMode: &accessMode,
	}
}

func DeveloperToPublicListing(dev ApiDeveloperListing) ApiPublicListing {
	return ApiPublicListing{ApiListingBase: dev.ApiListingBase}
}

// --------------------
// Handlers
// --------------------

func (api *DeveloperListingAPI) ListListingsHandler(userID shared.UserID, c echo.Context) error {
	return util.HandleBatchRequest(
		c,
		func(
			ctx context.Context,
			request shared.BatchRequest[listing.ListingCursor],
		) (items []*listing.Listing, nextCursor listing.ListingCursor, err error) {
			return api.listingService.FetchNextBatchByAuthor(ctx, userID, request)
		},
		func(elements []*listing.Listing) (views []ApiDeveloperListing) {
			return util.MapList(elements, MakeApiDeveloperListing)
		},
	)
}

func (api *DeveloperListingAPI) GetListingHandler(userID shared.UserID, c echo.Context) error {
	ctx := c.Request().Context()
	listingID := c.Param("id")
	if listingID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Listing ID is required")
	}

	l, err := api.listingService.GetOwnedListing(ctx, shared.ListingID(listingID), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Listing not found")
	}

	return c.JSON(http.StatusOK, MakeApiDeveloperListing(l))
}

func (api *DeveloperListingAPI) AddListingHandler(userID shared.UserID, c echo.Context) error {
	ctx := c.Request().Context()
	var req AddListingRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	newListing, err := listing.NewListing(userID, req.Title, req.Description, listing.Private, &shared.HardwareSpecification{}, 0)

	if err = api.listingService.CreateListing(ctx, newListing); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create listing")
	}

	return c.JSON(http.StatusOK, MakeApiDeveloperListing(newListing))
}

func (api *DeveloperListingAPI) UpdateListingHandler(userID shared.UserID, c echo.Context) error {
	ctx := c.Request().Context()
	var req UpdateListingRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	updated, err := api.listingService.UpdateListing(ctx, req.ID, func(l *listing.Listing) error {
		if req.Title != nil {
			l.SetTitle(*req.Title)
		}
		if req.Description != nil {
			l.SetDescription(*req.Description)
		}
		if req.Hardware != nil {
			l.SetHardware(req.Hardware)
		}
		if req.Price != nil {
			l.SetPriceInMinorUnits(*req.Price)
		}
		if req.AccessMode != nil {
			l.SetAccessMode(*req.AccessMode)
		}

		return nil
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update listing")
	}

	return c.JSON(http.StatusOK, MakeApiDeveloperListing(updated))
}

func (api *DeveloperListingAPI) DeleteListingHandler(userID shared.UserID, c echo.Context) error {
	ctx := c.Request().Context()
	listingID := c.Param("id")
	if listingID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Listing ID is required")
	}

	if err := api.listingService.DeleteListing(ctx, shared.ListingID(listingID)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not delete listing")
	}

	return c.NoContent(http.StatusNoContent)
}
