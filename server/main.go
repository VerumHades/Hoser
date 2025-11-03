package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"server/internal/configuration"
	"server/internal/database"
	"server/internal/utils"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var runningConfiguration = configuration.Load()

var store = sessions.NewCookieStore([]byte(runningConfiguration.SessionSecret))
var databaseInteractor database.Interactor = &database.DummyInteractor{}

func sendIvalidCredentialsError(w http.ResponseWriter) {
	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	user, user_error := databaseInteractor.GetUserByName(username)

	if user_error != nil {
		sendIvalidCredentialsError(w)
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.GetPasswordHash()), []byte(password))

	if err != nil {
		sendIvalidCredentialsError(w)
		return
	}

	session, _ := store.Get(r, "user-session")

	session.Values["authenticated"] = true
	session.Values["user"] = user
	session.Save(r, w)

	//fmt.Println("Authentificated user: " + user.Username)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	session, _ := store.Get(r, "user-session")

	session.Values["authenticated"] = false
	session.Options.MaxAge = -1
	session.Save(r, w)
}

func userDataRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}

	session, _ := store.Get(r, "user-session")
	//fmt.Print(session.Values)
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

type DeveloperListings struct {
	Author      string
	Title       string
	Description string
}

func developerListingsRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}

	session, _ := store.Get(r, "user-session")
	user, _ := session.Values["user"].(database.User)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	db_listings, db_error := user.GetListings()

	if db_error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	listings := make([]DeveloperListings, len(db_listings))

	username := user.GetUsername()
	for i := range db_listings {
		listings[i] = DeveloperListings{
			Author:      username,
			Description: db_listings[i].GetDescription(),
			Title:       db_listings[i].GetTitle(),
		}
	}

	json.NewEncoder(w).Encode(listings)
}

func developerAddListingRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	type AddListingRequest struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	var req AddListingRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	title := req.Title
	description := req.Description

	session, _ := store.Get(r, "user-session")
	user, _ := session.Values["user"].(database.User)

	listing, err := user.CreateListing(title, description)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(DeveloperListings{
		Author:      user.GetUsername(),
		Description: listing.GetDescription(),
		Title:       listing.GetTitle(),
	})
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

func authRequired(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "user-session")
		//fmt.Print(session.Values)
		// Check if authenticated value exists and is true
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Continue to next handler
		next.ServeHTTP(w, r)
	})
}

// After auth check for developer
func developerOnly(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "user-session")
		user, _ := session.Values["user"].(database.User)

		if !user.IsDeveloper() {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Continue to next handler
		next.ServeHTTP(w, r)
	})
}

func withCORS(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		for _, o := range runningConfiguration.AllowedOrigins {
			if origin == o {
				w.Header().Set("Access-Control-Allow-Origin", origin) // echo back single origin
				w.Header().Set("Access-Control-Allow-Credentials", "true")
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

type Middleware = func(http.Handler) http.HandlerFunc

func applyMiddlewares(middlewares []Middleware, handlers map[string]http.HandlerFunc) {
	for i := range middlewares {
		for key, value := range handlers {
			handlers[key] = middlewares[i](http.HandlerFunc(value))
		}
	}
}

func main() {
	public_handlers := map[string]http.HandlerFunc{
		"/login":          loginHandler,
		"/logout":         logoutHandler,
		"/rentals/public": publicRentalQueryHandler,
	}

	login_required_handlers := map[string]http.HandlerFunc{
		"/user/data": userDataRequestHandler,
	}

	developer_only_handlers := map[string]http.HandlerFunc{
		"/developer/listings": developerListingsRequestHandler,
		"/developer/listing":  developerAddListingRequestHandler,
	}

	store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   true,                  // must use HTTPS in production
		SameSite: http.SameSiteNoneMode, // allow cross-origin requests
	}

	gob.Register(&database.DummyUser{})

	applyMiddlewares([]Middleware{developerOnly, authRequired}, developer_only_handlers)
	applyMiddlewares([]Middleware{authRequired}, login_required_handlers)

	handlers := utils.MergeMultipleMaps([]map[string]http.HandlerFunc{public_handlers, login_required_handlers, developer_only_handlers})

	applyMiddlewares([]Middleware{utils.ColorLogMiddleware, withCORS}, handlers)

	mux := http.NewServeMux()

	for key, handler := range handlers {
		mux.Handle(key, http.HandlerFunc(handler))
	}

	fs := http.FileServer(http.Dir(runningConfiguration.ClientDirectory))
	mux.Handle("/", reactHandler(fs))

	var address = fmt.Sprintf("%s:%s", runningConfiguration.Address, runningConfiguration.Port)
	fmt.Printf("Server running on http://%s\n", address)
	http.ListenAndServe(address, mux)
}
