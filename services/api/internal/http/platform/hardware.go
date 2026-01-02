package platformapi

import (
	"api/internal/http/authentification"
	"api/pkg/util"
	"context"
	"net/http"
	"time"

	"common/pkg/domain/rates"
	"common/pkg/shared"

	"github.com/labstack/echo/v4"
)

// --------------------
// Repository interface
// --------------------

type APIHardwareCostQueryRepository interface {
	GetActiveRate(
		ctx context.Context,
		resourceType rates.HardwareResourceType,
		at time.Time,
	) (*rates.HardwareCostRate, error)

	FetchNextBatchOrderedByEffectiveDate(
		ctx context.Context,
		request shared.BatchRequest[rates.HardwareCostRateCursor],
	) (rates []*rates.HardwareCostRate, nextCursor rates.HardwareCostRateCursor, err error)
}

// --------------------
// API
// --------------------

type HardwareCostRatesAPI struct {
	queryRepository APIHardwareCostQueryRepository
}

func NewHardwareCostRatesAPI(
	queryRepository APIHardwareCostQueryRepository,
) *HardwareCostRatesAPI {
	return &HardwareCostRatesAPI{
		queryRepository: queryRepository,
	}
}

// --------------------
// Route registration
// --------------------

func (api *HardwareCostRatesAPI) RegisterRoutes(group *echo.Group) {
	group.GET(
		"/hardware-costs/active",
		authentification.WithAuthenticatedUser(api.GetActiveRateHandler),
	)
	group.GET(
		"/hardware-costs",
		authentification.WithAuthenticatedUser(api.ListRatesHandler),
	)
}

// --------------------
// DTOs
// --------------------

type HardwareCostRateResponse struct {
	ID            shared.HardwareCostRateID  `json:"id"`
	ResourceType  rates.HardwareResourceType `json:"resourceType"`
	CostInCents   int64                      `json:"costInCents"`
	ValidFromTime time.Time                  `json:"validFrom"`
}

// --------------------
// Conversion helpers
// --------------------

func MakeHardwareCostRateResponse(rate *rates.HardwareCostRate) HardwareCostRateResponse {
	return HardwareCostRateResponse{
		ID:            rate.ID(),
		ResourceType:  rate.Resource(),
		CostInCents:   rate.CostInCents(),
		ValidFromTime: rate.ValidFrom(),
	}
}

// --------------------
// Handlers
// --------------------

func (api *HardwareCostRatesAPI) GetActiveRateHandler(
	_ shared.UserID,
	c echo.Context,
) error {
	ctx := c.Request().Context()

	resourceType := rates.HardwareResourceType(c.QueryParam("resourceType"))
	if resourceType == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "resourceType is required")
	}

	atTime := time.Now()
	if atParam := c.QueryParam("at"); atParam != "" {
		parsedTime, err := time.Parse(time.RFC3339, atParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid at timestamp")
		}
		atTime = parsedTime
	}

	rate, err := api.queryRepository.GetActiveRate(ctx, resourceType, atTime)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "active rate not found")
	}

	return c.JSON(
		http.StatusOK,
		MakeHardwareCostRateResponse(rate),
	)
}

func (api *HardwareCostRatesAPI) ListRatesHandler(
	_ shared.UserID,
	c echo.Context,
) error {
	return util.HandleBatchRequest(
		c,
		func(
			ctx context.Context,
			request shared.BatchRequest[rates.HardwareCostRateCursor],
		) (rates []*rates.HardwareCostRate, nextCursor rates.HardwareCostRateCursor, err error) {
			return api.queryRepository.FetchNextBatchOrderedByEffectiveDate(ctx, request)
		},
		func(elements []*rates.HardwareCostRate) []HardwareCostRateResponse {
			return util.MapList(elements, MakeHardwareCostRateResponse)
		},
	)
}
