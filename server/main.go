package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"server/internal/configuration"
	"server/internal/database"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte("super-secret-key"))
var databaseInteractor database.Interactor = &database.DummyInteractor{}
var runningConfiguration = configuration.Load()

func sendIvalidCredentialsError(w http.ResponseWriter) {
	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	user, user_error := databaseInteractor.GetUserByName(username)

	if user_error != nil {
		sendIvalidCredentialsError(w)
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))

	if err != nil {
		sendIvalidCredentialsError(w)
		return
	}

	session, _ := store.Get(r, "user-session")

	session.Values["authenticated"] = true
	session.Values["user"] = user
	session.Save(r, w)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	session, _ := store.Get(r, "user-session")

	session.Values["authenticated"] = false
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func publicRentalQueryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}

	options := database.RentalQueryOptions{}

	query := r.URL.Query() // returns url.Values (map[string][]string)

	if q := query.Get("q"); q != "" {
		options.Text = q
	}

	response, err := databaseInteractor.QueryPublicRentals(&options)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

func withCORS(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		for _, o := range runningConfiguration.AllowedOrigins {
			if origin == o {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func reactHandler(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(runningConfiguration.ClientDirectory, r.URL.Path)

		if _, err := os.Stat(path); os.IsNotExist(err) {
			http.ServeFile(w, r, filepath.Join(runningConfiguration.ClientDirectory, "index.html"))
			return
		}

		fs.ServeHTTP(w, r)
	}
}

func main() {
	handlers := map[string]func(http.ResponseWriter, *http.Request){
		"/login":          loginHandler,
		"/logout":         logoutHandler,
		"/rentals/public": publicRentalQueryHandler,
	}
	middlewares := []func(http.Handler) http.HandlerFunc{
		withCORS,
	}

	mux := http.NewServeMux()

	for key, value := range handlers {
		handler := value

		for i := range middlewares {
			handler = middlewares[i](http.HandlerFunc(handler))
		}

		mux.Handle(key, http.HandlerFunc(handler))
	}

	fs := http.FileServer(http.Dir(runningConfiguration.ClientDirectory))
	mux.Handle("/", reactHandler(fs))

	var address = fmt.Sprintf("%s:%s", runningConfiguration.Address, runningConfiguration.Port)
	fmt.Printf("Server running on http://%s\n", address)
	http.ListenAndServe(address, mux)
}
