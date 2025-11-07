package handlers

import (
	"net/http"
	"server/internal/configuration"
	"server/internal/database"

	"github.com/gorilla/sessions"
)

type App struct {
	Store                *sessions.CookieStore
	DatabaseInteractor   database.Interactor
	RunningConfiguration *configuration.Configuration
}

func NewApp(store *sessions.CookieStore, db database.Interactor, cfg *configuration.Configuration) *App {
	var app = &App{
		Store:                store,
		DatabaseInteractor:   db,
		RunningConfiguration: cfg,
	}

	app.Store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   true,                  // must use HTTPS in production
		SameSite: http.SameSiteNoneMode, // allow cross-origin requests
	}

	return app
}
