package developerapi

import (
	"api/internal/http/authentification"
	"context"
	"net/http"

	"common/pkg/shared"

	"github.com/labstack/echo/v4"
)

// --------------------
// Services
// --------------------

type APIDeveloperListingSetupService interface {
	GetPrivateSetupForOwnedListing(
		ctx context.Context,
		userID shared.UserID,
		listingID shared.ListingID,
	) (*ListingGithubSetup, error)

	AttachOrUpdateSetup(
		ctx context.Context,
		userID shared.UserID,
		listingID shared.ListingID,
		repositoryURL string,
		accessToken string,
	) (*ListingGithubSetup, error)

	RemoveSetup(
		ctx context.Context,
		userID shared.UserID,
		listingID shared.ListingID,
	) error
}

// --------------------
// API
// --------------------

type DeveloperListingSetupAPI struct {
	listingSetupService APIDeveloperListingSetupService
}

func NewDeveloperListingSetupAPI(
	listingSetupService APIDeveloperListingSetupService,
) *DeveloperListingSetupAPI {
	return &DeveloperListingSetupAPI{
		listingSetupService: listingSetupService,
	}
}

// --------------------
// Route registration
// --------------------

func (api *DeveloperListingSetupAPI) RegisterRoutes(group *echo.Group) {
	group.GET(
		"/listings/:id/setup/github",
		authentification.WithAuthenticatedUser(api.GetGithubSetupHandler),
	)
	group.PUT(
		"/listings/:id/setup/github",
		authentification.WithAuthenticatedUser(api.AttachOrUpdateGithubSetupHandler),
	)
	group.DELETE(
		"/listings/:id/setup/github",
		authentification.WithAuthenticatedUser(api.RemoveGithubSetupHandler),
	)
}

// --------------------
// DTOs
// --------------------

type GithubSetupResponse struct {
	RepositoryURL string `json:"repoUrl"`
}

type AttachOrUpdateGithubSetupRequest struct {
	RepositoryURL string `json:"repoUrl"`
	AccessToken   string `json:"accessToken"`
}

// --------------------
// Domain projection
// --------------------

type ListingGithubSetup struct {
	RepositoryURL string
}

// --------------------
// Handlers
// --------------------

func (api *DeveloperListingSetupAPI) GetGithubSetupHandler(
	userID shared.UserID,
	c echo.Context,
) error {
	ctx := c.Request().Context()
	listingID := shared.ListingID(c.Param("id"))
	if listingID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Listing ID is required")
	}

	setup, err := api.listingSetupService.GetPrivateSetupForOwnedListing(
		ctx,
		userID,
		listingID,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "GitHub setup not found")
	}

	return c.JSON(
		http.StatusOK,
		GithubSetupResponse{
			RepositoryURL: setup.RepositoryURL,
		},
	)
}

func (api *DeveloperListingSetupAPI) AttachOrUpdateGithubSetupHandler(
	userID shared.UserID,
	c echo.Context,
) error {
	ctx := c.Request().Context()
	listingID := shared.ListingID(c.Param("id"))
	if listingID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Listing ID is required")
	}

	var request AttachOrUpdateGithubSetupRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	setup, err := api.listingSetupService.AttachOrUpdateSetup(
		ctx,
		userID,
		listingID,
		request.RepositoryURL,
		request.AccessToken,
	)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Failed to attach or update GitHub setup",
		)
	}

	return c.JSON(
		http.StatusOK,
		GithubSetupResponse{
			RepositoryURL: setup.RepositoryURL,
		},
	)
}

func (api *DeveloperListingSetupAPI) RemoveGithubSetupHandler(
	userID shared.UserID,
	c echo.Context,
) error {
	ctx := c.Request().Context()
	listingID := shared.ListingID(c.Param("id"))
	if listingID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Listing ID is required")
	}

	if err := api.listingSetupService.RemoveSetup(ctx, userID, listingID); err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Failed to remove GitHub setup",
		)
	}

	return c.NoContent(http.StatusNoContent)
}
