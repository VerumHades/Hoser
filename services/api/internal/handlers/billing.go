package handlers

import (
	"net/http"

	"common/pkg/billing/account"

	"github.com/labstack/echo/v4"
)

// =================== TYPES ===================
type ApiBillingAccount struct {
	ID                string `json:"id"`
	Status            string `json:"status"`
	PaymentProvider   string `json:"paymentProvider"`
	ProviderAccountID string `json:"providerAccountId"`
	CreatedAt         int64  `json:"createdAt"` // unix timestamp
}

type CreateBillingAccountRequest struct {
	PaymentProvider   string `json:"paymentProvider"`
	ProviderAccountID string `json:"providerAccountId"`
}

// =================== HELPERS ===================
func (app *App) MakeApiBillingAccount(account *account.BillingAccount) ApiBillingAccount {
	return ApiBillingAccount{
		ID:                account.ID,
		Status:            string(account.Status),
		PaymentProvider:   string(account.PaymentProvider),
		ProviderAccountID: account.ProviderAccountID,
		CreatedAt:         account.CreatedAt.Unix(),
	}
}

// =================== HANDLERS ===================

// UserBillingAccountsHandler returns all billing accounts for the authenticated user.
func (app *App) UserBillingAccountsHandler(c echo.Context) error {
	userID, err := app.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	accounts, err := app.BillingAccountService.ListByOwner(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch billing accounts")
	}

	apiAccounts := make([]ApiBillingAccount, len(accounts))
	for i, acct := range accounts {
		apiAccounts[i] = app.MakeApiBillingAccount(acct)
	}

	return c.JSON(http.StatusOK, apiAccounts)
}

// UserCreateBillingAccountHandler creates a new billing account for the authenticated user.
func (app *App) UserCreateBillingAccountHandler(c echo.Context) error {
	userID, err := app.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	var req CreateBillingAccountRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	newAccount, err := app.BillingAccountService.CreateAccount(
		userID,
		"",
		req.ProviderAccountID,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create billing account")
	}

	return c.JSON(http.StatusOK, app.MakeApiBillingAccount(newAccount))
}

// UserGetBillingAccountHandler returns a single billing account for the authenticated user.
func (app *App) UserGetBillingAccountHandler(c echo.Context) error {
	userID, err := app.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	accountID := c.Param("id")
	acct, err := app.BillingAccountService.GetAccount(accountID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch billing account")
	}
	if acct == nil || acct.OwnerID != userID {
		return echo.NewHTTPError(http.StatusNotFound, "Billing account not found")
	}

	return c.JSON(http.StatusOK, app.MakeApiBillingAccount(acct))
}

// UserSuspendBillingAccountHandler suspends a billing account for the authenticated user.
func (app *App) UserSuspendBillingAccountHandler(c echo.Context) error {
	userID, err := app.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	accountID := c.Param("id")
	acct, err := app.BillingAccountService.GetAccount(accountID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch billing account")
	}
	if acct == nil || acct.OwnerID != userID {
		return echo.NewHTTPError(http.StatusNotFound, "Billing account not found")
	}

	if err := app.BillingAccountService.SuspendAccount(accountID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to suspend billing account")
	}

	return c.JSON(http.StatusOK, app.MakeApiBillingAccount(acct))
}

// UserCloseBillingAccountHandler closes a billing account for the authenticated user.
func (app *App) UserCloseBillingAccountHandler(c echo.Context) error {
	userID, err := app.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	accountID := c.Param("id")
	acct, err := app.BillingAccountService.GetAccount(accountID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch billing account")
	}
	if acct == nil || acct.OwnerID != userID {
		return echo.NewHTTPError(http.StatusNotFound, "Billing account not found")
	}

	if err := app.BillingAccountService.CloseAccount(accountID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to close billing account")
	}

	return c.JSON(http.StatusOK, app.MakeApiBillingAccount(acct))
}
