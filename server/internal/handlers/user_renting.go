package handlers

import (
	"encoding/json"
	"net/http"
)

type ApiUserRental struct {
	ID              string
	Title           string
	Description     string
	SourceListingID string
}

func (app *App) UserRentalsRequestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := app.GetUser(r)

		dbRentals, dbErr := user.GetRentals()
		if dbErr != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		rentals := make([]ApiUserRental, len(dbRentals))
		for i, rental := range dbRentals {
			rentals[i] = ApiUserRental{
				ID:              rental.GetUUID(),
				Title:           rental.GetTitle(),
				Description:     rental.GetDescription(),
				SourceListingID: rental.GetSourceListingUUID(),
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rentals)
	}
}

func (app *App) UserRentRequestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type RentalRequest struct {
			ID string `json:"id"` // required
		}

		var req RentalRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if req.ID == "" {
			http.Error(w, "Missing id", http.StatusBadRequest)
			return
		}

		_, dbErr := app.DatabaseInteractor.GetPublicListing(req.ID)
		if dbErr != nil {
			http.Error(w, "Listing not found", http.StatusNotFound)
			return
		}

		user := app.GetUser(r)
		user.Rent(req.ID)

		w.WriteHeader(http.StatusCreated)
	}
}
