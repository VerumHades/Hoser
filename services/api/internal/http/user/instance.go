package userapi

import (
	"api/internal/http/authentification"
	"context"
	"net/http"
	"time"

	"common/pkg/domain/instance"
	"common/pkg/shared"

	"github.com/labstack/echo/v4"
)

// --------------------
// Service interface
// --------------------

type APIUserInstanceContractService interface {
	RentInstanceOfListing(
		ctx context.Context,
		userID shared.UserID,
		listingID shared.ListingID,
		contract *instance.InstanceRentalContract,
	) error

	EnableContractRenewal(
		ctx context.Context,
		contractID shared.InstanceRentalContractID,
		renewDuration time.Duration,
	) error

	DisableContractRenewal(
		ctx context.Context,
		contractID shared.InstanceRentalContractID,
	) error

	ChangeContractHardwareSpecification(
		ctx context.Context,
		contractID shared.InstanceRentalContractID,
		hardwareSpecification *shared.HardwareSpecification,
	) error
}

// --------------------
// API
// --------------------

type InstanceContractAPI struct {
	contractService APIUserInstanceContractService
}

func NewInstanceContractAPI(
	contractService APIUserInstanceContractService,
) *InstanceContractAPI {
	return &InstanceContractAPI{
		contractService: contractService,
	}
}

// --------------------
// Route registration
// --------------------

func (api *InstanceContractAPI) RegisterRoutes(group *echo.Group) {
	group.POST(
		"/contracts",
		authentification.WithAuthenticatedUser(api.CreateContractHandler),
	)
	group.PUT(
		"/contracts/:id/renewal/enable",
		authentification.WithAuthenticatedUser(api.EnableRenewalHandler),
	)
	group.PUT(
		"/contracts/:id/renewal/disable",
		authentification.WithAuthenticatedUser(api.DisableRenewalHandler),
	)
	group.PUT(
		"/contracts/:id/hardware",
		authentification.WithAuthenticatedUser(api.ChangeHardwareHandler),
	)
}

// --------------------
// DTOs
// --------------------

type CreateContractRequest struct {
	ListingID             shared.ListingID              `json:"listingId"`
	HardwareSpecification *shared.HardwareSpecification `json:"hardware"`
	Duration              time.Duration                 `json:"duration"`
	RenewAutomatically    bool                          `json:"renewAutomatically"`
	RenewalDuration       time.Duration                 `json:"renewalDuration"`
}
type EnableRenewalRequest struct {
	RenewalDuration time.Duration `json:"renewalDuration"`
}

type ChangeHardwareRequest struct {
	HardwareSpecification *shared.HardwareSpecification `json:"hardware"`
}

// --------------------
// Handlers
// --------------------
func (api *InstanceContractAPI) CreateContractHandler(
	userID shared.UserID,
	c echo.Context,
) error {
	ctx := c.Request().Context()

	var request CreateContractRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	if request.Duration <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "duration must be positive")
	}

	if request.HardwareSpecification == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "hardware specification is required")
	}

	periodStart := time.Now().UTC()
	periodEnd := periodStart.Add(request.Duration)

	contract, err := instance.NewInstanceRentalContract(
		request.ListingID,
		userID,
		request.HardwareSpecification,
		periodStart,
		periodEnd,
		request.RenewAutomatically,
		request.RenewalDuration,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := api.contractService.RentInstanceOfListing(
		ctx,
		userID,
		request.ListingID,
		contract,
	); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusCreated)
}

func (api *InstanceContractAPI) EnableRenewalHandler(
	_ shared.UserID,
	c echo.Context,
) error {
	ctx := c.Request().Context()
	contractID := shared.InstanceRentalContractID(c.Param("id"))
	if contractID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Contract ID is required")
	}

	var request EnableRenewalRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	if err := api.contractService.EnableContractRenewal(
		ctx,
		contractID,
		request.RenewalDuration,
	); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (api *InstanceContractAPI) DisableRenewalHandler(
	_ shared.UserID,
	c echo.Context,
) error {
	ctx := c.Request().Context()
	contractID := shared.InstanceRentalContractID(c.Param("id"))
	if contractID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Contract ID is required")
	}

	if err := api.contractService.DisableContractRenewal(ctx, contractID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (api *InstanceContractAPI) ChangeHardwareHandler(
	_ shared.UserID,
	c echo.Context,
) error {
	ctx := c.Request().Context()
	contractID := shared.InstanceRentalContractID(c.Param("id"))
	if contractID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Contract ID is required")
	}

	var request ChangeHardwareRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	if request.HardwareSpecification == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "hardware specification is required")
	}

	if err := api.contractService.ChangeContractHardwareSpecification(
		ctx,
		contractID,
		request.HardwareSpecification,
	); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
