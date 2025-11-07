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

func (app *App) GetUser(request *http.Request) database.User {
	session, _ := app.Store.Get(request, "user-session")
	user, _ := session.Values["user"].(database.User)
	return user
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

func NewTestApp() *App {
	var store = sessions.NewCookieStore([]byte("secret-key"))
	var db = &database.DummyInteractor{}
	var cfg = &configuration.Configuration{
		AllowedOrigins: []string{"http://localhost"},
	}

	return NewApp(store, db, cfg)
}
