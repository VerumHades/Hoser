package handlers

import (
	"net/http"

	"common/pkg/billing/payments/payment"

	"github.com/labstack/echo/v4"
)

// =================== TYPES ===================
type ApiPayment struct {
	ID               string  `json:"id"`
	BillingAccountID string  `json:"billingAccountId"`
	Amount           float64 `json:"amount"`
	Currency         string  `json:"currency"`
	Status           string  `json:"status"`
	Kind             string  `json:"kind"`
	CreatedAt        int64   `json:"createdAt"` // unix timestamp
	PaidAt           *int64  `json:"paidAt,omitempty"`
}

// =================== HELPERS ===================
func (app *App) MakeApiPayment(payment *payment.Payment) ApiPayment {
	var paidAt *int64
	if payment.PaidAt != nil {
		unix := payment.PaidAt.Unix()
		paidAt = &unix
	}

	return ApiPayment{
		ID:               payment.ID,
		BillingAccountID: payment.BillingAccountID,
		Amount:           payment.Amount.Amount,
		Currency:         payment.Amount.CurrencyCode,
		Status:           string(payment.Status),
		Kind:             string(payment.Kind),
		CreatedAt:        payment.CreatedAt.Unix(),
		PaidAt:           paidAt,
	}
}

// =================== HANDLERS ===================

// List all payments for a specific billing account
func (app *App) UserListPaymentsHandler(c echo.Context) error {
	userID, err := app.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	accountID := c.Param("billingAccountId")
	acct, err := app.BillingAccountService.GetAccount(accountID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch billing account")
	}
	if acct == nil || acct.OwnerID != userID {
		return echo.NewHTTPError(http.StatusNotFound, "Billing account not found")
	}

	payments, err := app.PaymentService.ListByBillingAccount(accountID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch payments")
	}

	apiPayments := make([]ApiPayment, len(payments))
	for i, p := range payments {
		apiPayments[i] = app.MakeApiPayment(p)
	}

	return c.JSON(http.StatusOK, apiPayments)
}

// Get a single payment
func (app *App) UserGetPaymentHandler(c echo.Context) error {
	userID, err := app.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	paymentID := c.Param("paymentId")
	paymentEntity, err := app.PaymentService.GetPaymentView(paymentID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch payment")
	}
	if paymentEntity == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Payment not found")
	}

	acct, err := app.BillingAccountService.GetAccount(paymentEntity.BillingAccountID)
	if err != nil || acct == nil || acct.OwnerID != userID {
		return echo.NewHTTPError(http.StatusNotFound, "Payment not found")
	}

	return c.JSON(http.StatusOK, app.MakeApiPayment(paymentEntity))
}

// Optionally, endpoint for metadata
func (app *App) UserGetPaymentMetadataHandler(c echo.Context) error {
	userID, err := app.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	paymentID := c.Param("paymentId")
	paymentEntity, err := app.PaymentService.GetPaymentView(paymentID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch payment")
	}
	if paymentEntity == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Payment not found")
	}

	acct, err := app.BillingAccountService.GetAccount(paymentEntity.BillingAccountID)
	if err != nil || acct == nil || acct.OwnerID != userID {
		return echo.NewHTTPError(http.StatusNotFound, "Payment not found")
	}

	switch paymentEntity.Kind {
	case payment.PaymentTypeOneTime:
		meta, err := app.PaymentService.GetOneTimePaymentMetadata(paymentID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch one-time metadata")
		}
		return c.JSON(http.StatusOK, meta)
	case payment.PaymentTypeSubscription:
		meta, err := app.PaymentService.GetSubscriptionPaymentMetadata(paymentID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch subscription metadata")
		}
		return c.JSON(http.StatusOK, meta)
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "Payment has no metadata")
	}
}
