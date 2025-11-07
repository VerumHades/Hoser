package handlers

import (
	"encoding/json"
	"net/http"
	"server/internal/database"
)

type ApiDeveloperListing struct {
	Author      string
	ID          string
	Title       string
	Description string
	AccessMode  database.ListingAccessMode
}

func makeApiDeveloperListing(authorName string, listing database.Listing) ApiDeveloperListing {
	return ApiDeveloperListing{
		Author:      authorName,
		ID:          listing.GetUUID(),
		Description: listing.GetDescription(),
		Title:       listing.GetTitle(),
		AccessMode:  listing.GetAccessMode(),
	}
}

func (app *App) DeveloperListingsRequestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := app.Store.Get(r, "user-session")
		user, _ := session.Values["user"].(database.User)

		dbListings, dbErr := user.GetListings()
		if dbErr != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		listings := make([]ApiDeveloperListing, len(dbListings))
		username := user.GetUsername()
		for i, l := range dbListings {
			listings[i] = makeApiDeveloperListing(username, l)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(listings)
	}
}

func (app *App) DeveloperAddListingRequestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type AddListingRequest struct {
			Title       string `json:"title"`
			Description string `json:"description"`
		}

		var req AddListingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		session, _ := app.Store.Get(r, "user-session")
		user, _ := session.Values["user"].(database.User)

		listing, err := user.CreateListing(req.Title, req.Description)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(makeApiDeveloperListing(user.GetUsername(), listing))
	}
}

type ListingRequest struct {
	ID          string  `json:"id"`                    // required
	Title       *string `json:"title,omitempty"`       // optional
	Description *string `json:"description,omitempty"` // optional
	AccessMode  *int    `json:"accessMode,omitempty"`
}

func (app *App) DeveloperAlterListingRequestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ListingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if req.ID == "" {
			http.Error(w, "Missing listing ID", http.StatusBadRequest)
			return
		}

		session, _ := app.Store.Get(r, "user-session")
		user, _ := session.Values["user"].(database.User)

		listing, err := user.GetListing(req.ID)
		if err != nil {
			http.Error(w, "Listing not found", http.StatusNotFound)
			return
		}

		if req.Title != nil {
			listing.SetTitle(*req.Title)
		}
		if req.Description != nil {
			listing.SetDescription(*req.Description)
		}
		if req.AccessMode != nil {
			mode, err := database.ListingAccessModeFromInt(*req.AccessMode)
			if err != nil {
				http.Error(w, "Invalid access mode.", http.StatusBadRequest)
				return
			}
			listing.SetAccessMode(mode)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(makeApiDeveloperListing(user.GetUsername(), listing))
	}
}

func (app *App) DeveloperDeleteListingRequestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ListingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if req.ID == "" {
			http.Error(w, "Missing listing ID", http.StatusBadRequest)
			return
		}

		session, _ := app.Store.Get(r, "user-session")
		user, _ := session.Values["user"].(database.User)

		user.DeleteListing(req.ID)

		w.WriteHeader(http.StatusNoContent)
	}
}
