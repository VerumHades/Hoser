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
	"server/internal/handlers"
	"server/internal/utils"

	"github.com/gorilla/sessions"
)

var runningConfiguration = configuration.Load()

var app = handlers.App{
	RunningConfiguration: &runningConfiguration,
	Store:                sessions.NewCookieStore([]byte(runningConfiguration.SessionSecret)),
	DatabaseInteractor:   &database.DummyInteractor{},
}

func publicRentalQueryHandler(w http.ResponseWriter, r *http.Request) {
	options := database.RentalQueryOptions{}

	query := r.URL.Query() // returns url.Values (map[string][]string)

	if q := query.Get("q"); q != "" {
		options.Text = q
	}

	response, err := app.DatabaseInteractor.QueryPublicRentals(&options)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
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

type HandlerMethodMap = map[string]http.HandlerFunc
type HandlerMap = map[string]HandlerMethodMap

func buildHandlerForAllRequestMethods(method_map HandlerMethodMap) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler, exists := method_map[r.Method]

		if !exists {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func buildHandlerMap(handler_map HandlerMap) map[string]http.HandlerFunc {
	built_map := map[string]http.HandlerFunc{}

	for route, method_map := range handler_map {
		built_map[route] = buildHandlerForAllRequestMethods(method_map)
	}

	return built_map
}

func main() {
	public_handlers := buildHandlerMap(HandlerMap{
		"/login":          {http.MethodPost: app.LoginHandler()},
		"/logout":         {http.MethodPost: app.LogoutHandler()},
		"/rentals/public": {http.MethodGet: publicRentalQueryHandler},
	})

	login_required_handlers := buildHandlerMap(HandlerMap{
		"/user/data": {http.MethodGet: app.UserDataRequestHandler()},
	})

	developer_only_handlers := buildHandlerMap(HandlerMap{
		"/developer/listings": {http.MethodGet: app.DeveloperListingsRequestHandler()},
		"/developer/listing": {
			http.MethodPost:   app.DeveloperAddListingRequestHandler(),
			http.MethodPut:    app.DeveloperAlterListingRequestHandler(),
			http.MethodDelete: app.DeveloperDeleteListingRequestHandler(),
		},
	})

	gob.Register(&database.DummyUser{})

	handlers.ApplyMiddlewares([]handlers.Middleware{app.MiddlewareDeveloperOnly(), app.MiddlewareAuthenticationRequired()}, developer_only_handlers)
	handlers.ApplyMiddlewares([]handlers.Middleware{app.MiddlewareAuthenticationRequired()}, login_required_handlers)

	all_handlers := utils.MergeMultipleMaps([]map[string]http.HandlerFunc{public_handlers, login_required_handlers, developer_only_handlers})

	handlers.ApplyMiddlewares([]handlers.Middleware{utils.ColorLogMiddleware, app.MiddlewareCORS()}, all_handlers)

	mux := http.NewServeMux()

	for key, handler := range all_handlers {
		mux.Handle(key, http.HandlerFunc(handler))
	}

	fs := http.FileServer(http.Dir(runningConfiguration.ClientDirectory))
	mux.Handle("/", reactHandler(fs))

	var address = fmt.Sprintf("%s:%s", runningConfiguration.Address, runningConfiguration.Port)
	fmt.Printf("Server running on http://%s\n", address)
	http.ListenAndServe(address, mux)
}
