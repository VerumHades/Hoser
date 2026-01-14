package developerapi

import (
	"api/internal/http/authentification"
	"api/pkg/util"
	"common/pkg/application/services/developer"
	"common/pkg/domain/entities/listing"
	"common/pkg/domain/repositories"
	"common/pkg/shared"
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type APIDeveloperListingService interface {
	CreateListing(ctx context.Context, request developer.CreateListingRequest) (l *listing.Listing, err error)
	UpdateListing(ctx context.Context, listingID shared.ListingID, updateFunction developer.ListingMutationFunction) (l *listing.Listing, err error)
	DeleteListing(ctx context.Context, listingID shared.ListingID) error

	FetchNextBatchByAuthor(
		ctx context.Context,
		authorID shared.UserID,
		request shared.BatchRequest[repositories.ListingCursor],
	) (listings []*listing.Listing, nextCursor repositories.ListingCursor, err error)

	GetOwnedListing(ctx context.Context, listingID shared.ListingID, userID shared.UserID) (*listing.Listing, error)

	CreateScreenshotUploadUrl(ctx context.Context, listingID shared.ListingID, userID shared.UserID) (string, error)
	DeleteScreenshot(ctx context.Context, listingID shared.ListingID, userID shared.UserID, screenshotId shared.ListingScreenshotID) error
}

type APIDeveloperUserService interface {
	IsUserDeveloper(ctx context.Context, userID shared.UserID) (bool, error)
}

// DeveloperListingAPI is the API layer for developer-specific listing endpoints
type DeveloperListingAPI struct {
	listingService APIDeveloperListingService
	userService    APIDeveloperUserService
}

// NewDeveloperListingAPI constructs a DeveloperListingAPI with required dependencies.
func NewDeveloperListingAPI(
	listingService APIDeveloperListingService,
	userService APIDeveloperUserService,
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

	// Screenshot Routes
	group.POST("/listings/:id/screenshots", authentification.WithAuthenticatedUser(api.CreateScreenshotUploadHandler))
	group.DELETE("/listings/:id/screenshots/:screenshotId", authentification.WithAuthenticatedUser(api.DeleteScreenshotHandler))
}

// --------------------
// API DTOs
// --------------------

type ApiListingBase struct {
	ID            shared.ListingID              `json:"id"`
	Title         string                        `json:"title,omitempty"`
	Description   string                        `json:"description,omitempty"`
	Price         int64                         `json:"price"`
	Hardware      *shared.HardwareSpecification `json:"hardware,omitempty"`
	ScreenshotIds []shared.ListingScreenshotID  `json:"screenshotIds"`
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
			ID:            l.ID(),
			Title:         l.Title(),
			Description:   l.Description(),
			Price:         l.PriceInMinorUnits(),
			Hardware:      l.HardwareSpecification(),
			ScreenshotIds: l.ScreenshotKeys(),
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
			request shared.BatchRequest[repositories.ListingCursor],
		) (items []*listing.Listing, nextCursor repositories.ListingCursor, err error) {
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

	newListing, err := api.listingService.CreateListing(ctx, developer.CreateListingRequest{
		AuthorID:          userID,
		Title:             req.Title,
		Description:       req.Description,
		AccessMode:        listing.Private,
		Hardware:          &shared.HardwareSpecification{},
		PriceInMinorUnits: 1000,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create listing: "+err.Error())
	}

	return c.JSON(http.StatusOK, MakeApiDeveloperListing(newListing))
}

func (api *DeveloperListingAPI) UpdateListingHandler(userID shared.UserID, c echo.Context) error {
	ctx := c.Request().Context()
	var req UpdateListingRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	updated, err := api.listingService.UpdateListing(ctx, req.ID, func(ctx context.Context, mutator *listing.ListingMutationBuilder) error {
		if req.Title != nil {
			mutator = mutator.SetTitle(*req.Title)
		}
		if req.Description != nil {
			mutator = mutator.SetDescription(*req.Description)
		}
		if req.Hardware != nil {
			mutator = mutator.SetHardware(req.Hardware)
		}
		if req.Price != nil {
			mutator = mutator.SetPriceInMinorUnits(*req.Price)
		}
		if req.AccessMode != nil {
			mutator = mutator.SetAccessMode(*req.AccessMode)
		}

		return nil
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update listing: "+err.Error())
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

type CreateScreenshotResponse struct {
	UploadUrl string `json:"uploadUrl"`
}

func (api *DeveloperListingAPI) CreateScreenshotUploadHandler(userID shared.UserID, c echo.Context) error {
	ctx := c.Request().Context()
	listingID := c.Param("id")
	if listingID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Listing ID is required")
	}

	// Call service to validate ownership, check limits, and generate signed URL
	uploadUrl, err := api.listingService.CreateScreenshotUploadUrl(ctx, shared.ListingID(listingID), userID)
	if err != nil {
		if err == shared.ErrLimitReached {
			return echo.NewHTTPError(http.StatusConflict, "Maximum number of screenshots reached")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate upload URL")
	}

	return c.JSON(http.StatusCreated, CreateScreenshotResponse{
		UploadUrl: uploadUrl,
	})
}

func (api *DeveloperListingAPI) DeleteScreenshotHandler(userID shared.UserID, c echo.Context) error {
	ctx := c.Request().Context()
	listingID := c.Param("id")
	screenshotID := c.Param("screenshotId")

	if listingID == "" || screenshotID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Listing ID and Screenshot ID are required")
	}

	err := api.listingService.DeleteScreenshot(
		ctx,
		shared.ListingID(listingID),
		userID,
		shared.ListingScreenshotID(screenshotID),
	)

	if err != nil {
		// We assume 404/403 here since DeleteScreenshot verifies ownership
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not delete screenshot")
	}

	return c.NoContent(http.StatusNoContent)
}
