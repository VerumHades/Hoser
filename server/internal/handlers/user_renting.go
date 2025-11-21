package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// =================== TYPES ===================

type ApiUserRental struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	SourceListingID string `json:"sourceListingID"`
}

type RentalRequest struct {
	ID string `json:"id"`
}

// =================== HANDLERS ===================

// UserRentalsHandler returns all rentals of the authenticated user
func (app *App) UserRentalsHandler(c echo.Context) error {
	user, err := app.GetUserFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	dbRentals, dbErr := user.Rentals()
	if dbErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
	}

	rentals := make([]ApiUserRental, len(dbRentals))
	for i, rental := range dbRentals {
		rentals[i] = ApiUserRental{
			ID:              rental.UUID(),
			SourceListingID: rental.SourceListingUUID(),
		}
	}

	return c.JSON(http.StatusOK, rentals)
}

func (app *App) UserRentHandler(c echo.Context) error {
	user, err := app.GetUserFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	var req RentalRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	_, dbErr := app.DatabaseInteractor.GetPublicListing(req.ID)
	if dbErr != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Listing not found")
	}

	if err := user.Rent(req.ID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not rent listing")
	}

	return c.NoContent(http.StatusCreated)
}
