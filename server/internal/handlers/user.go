package handlers

import (
	"encoding/json"
	"net/http"
	"server/internal/database"

	"golang.org/x/crypto/bcrypt"
)

func sendInvalidCredentialsError(w http.ResponseWriter) {
	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
}

func (app *App) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")

		user, err := app.DatabaseInteractor.GetUserByName(username)
		if err != nil {
			sendInvalidCredentialsError(w)
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(user.GetPasswordHash()), []byte(password)) != nil {
			sendInvalidCredentialsError(w)
			return
		}

		session, _ := app.Store.Get(r, "user-session")
		session.Values["authenticated"] = true
		session.Values["user"] = user
		_ = session.Save(r, w)
	}
}

func (app *App) LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := app.Store.Get(r, "user-session")
		session.Values["authenticated"] = false

		session.Options.MaxAge = -1
		_ = session.Save(r, w)
	}
}

func (app *App) UserDataRequestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := app.Store.Get(r, "user-session")
		user, _ := session.Values["user"].(database.User)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(struct {
			Username    string
			IsDeveloper bool
		}{
			Username:    user.GetUsername(),
			IsDeveloper: user.IsDeveloper(),
		})
	}
}
