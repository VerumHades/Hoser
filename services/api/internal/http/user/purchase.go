package userapi

import (
	"api/internal/http/authentification"
	"context"
	"net/http"

	"common/pkg/shared"

	"github.com/labstack/echo/v4"
)

// --------------------
// Service interface
// --------------------

// APIUserPurchaseService defines the subset of functionality
// required by the API from the underlying domain services.
type APIUserPurchaseService interface {
	PurchaseListing(
		ctx context.Context,
		userID shared.UserID,
		listingID shared.ListingID,
	) error

	DoesUserOwnListing(
		ctx context.Context,
		userID shared.UserID,
		listingID shared.ListingID,
	) (bool, error)

	IsPurchaseProcessing(
		ctx context.Context,
		userID shared.UserID,
		listingID shared.ListingID,
	) (bool, error)
}

// --------------------
// API
// --------------------

type UserPurchaseAPI struct {
	purchaseService APIUserPurchaseService
}

func NewUserPurchaseAPI(
	purchaseService APIUserPurchaseService,
) *UserPurchaseAPI {
	return &UserPurchaseAPI{
		purchaseService: purchaseService,
	}
}

// --------------------
// Route registration
// --------------------

func (api *UserPurchaseAPI) RegisterRoutes(group *echo.Group) {
	// Purchase a listing
	group.POST(
		"/listings/:id/purchase",
		authentification.WithAuthenticatedUser(api.PurchaseListingHandler),
	)

	// Check if the current user owns a specific listing
	group.GET(
		"/listings/:id/ownership",
		authentification.WithAuthenticatedUser(api.CheckOwnershipHandler),
	)
}

// --------------------
// DTOs
// --------------------

type OwnershipResponse struct {
	IsOwner bool `json:"isOwner"`
}

type ProcessingResponse struct {
	IsProcessing bool `json:"isProcessing"`
}

// Note: PurchaseListing currently uses the ID from the URL,
// so a request DTO isn't strictly necessary unless you add options like payment methods.

// --------------------
// Handlers
// --------------------

func (api *UserPurchaseAPI) PurchaseListingHandler(
	userID shared.UserID,
	c echo.Context,
) error {
	ctx := c.Request().Context()
	listingID := shared.ListingID(c.Param("id"))
	if listingID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Listing ID is required")
	}

	if err := api.purchaseService.PurchaseListing(
		ctx,
		userID,
		listingID,
	); err != nil {
		// Depending on the error type from your service, you might want to
		// differentiate between 400 (Already owned) and 500 (Internal error)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusCreated)
}

func (api *UserPurchaseAPI) CheckOwnershipHandler(
	userID shared.UserID,
	c echo.Context,
) error {
	ctx := c.Request().Context()
	listingID := shared.ListingID(c.Param("id"))
	if listingID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Listing ID is required")
	}

	isOwner, err := api.purchaseService.DoesUserOwnListing(ctx, userID, listingID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, OwnershipResponse{
		IsOwner: isOwner,
	})
}

func (api *UserPurchaseAPI) IsPurchaseProcessingHandler(
	userID shared.UserID,
	c echo.Context,
) error {
	ctx := c.Request().Context()
	listingID := shared.ListingID(c.Param("id"))
	if listingID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Listing ID is required")
	}

	isProcessing, err := api.purchaseService.IsPurchaseProcessing(ctx, userID, listingID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, ProcessingResponse{
		IsProcessing: isProcessing,
	})
}
