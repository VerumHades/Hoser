package handlers

import (
	"common/pkg/rates"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type ApiMoney struct {
	Amount       float64 `json:"amount"`
	CurrencyCode string  `json:"currencyCode"`
}

type ApiHardwareRate struct {
	Amount       float64 `json:"amount"`
	CurrencyCode string  `json:"currencyCode"`
	Unit         string  `json:"unit"`
}

type ApiHardwareRatesResponse struct {
	CPUCost     ApiHardwareRate `json:"cpuCost"`
	RAMCost     ApiHardwareRate `json:"ramCost"`
	DiskCost    ApiHardwareRate `json:"diskCost"`
	EffectiveAt int64           `json:"effectiveAt"`
}

func MakeApiHardwareRates(rate *rates.HardwareCostRate) ApiHardwareRatesResponse {
	return ApiHardwareRatesResponse{
		CPUCost: ApiHardwareRate{
			Amount:       rate.CPUCost.Amount,
			CurrencyCode: rate.CPUCost.CurrencyCode,
			Unit:         "core-hour",
		},
		RAMCost: ApiHardwareRate{
			Amount:       rate.RAMCost.Amount,
			CurrencyCode: rate.RAMCost.CurrencyCode,
			Unit:         "byte-hour",
		},
		DiskCost: ApiHardwareRate{
			Amount:       rate.DiskCost.Amount,
			CurrencyCode: rate.DiskCost.CurrencyCode,
			Unit:         "byte-hour",
		},
		EffectiveAt: rate.ValidFrom.Unix(),
	}
}

func (app *App) HardwareRatesHandler(c echo.Context) error {
	currencyCode := c.QueryParam("currency")

	atParam := c.QueryParam("at")
	at := time.Now()

	if atParam != "" {
		parsed, err := strconv.ParseInt(atParam, 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid 'at' timestamp")
		}
		at = time.Unix(parsed, 0)
	}

	rateView, err := app.HardwareCostCalculationService.GetRateView(at, currencyCode)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch hardware rates")
	}

	return c.JSON(http.StatusOK, MakeApiHardwareRates(rateView))
}
