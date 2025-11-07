package handlers

import (
	"net/http"
	"server/internal/database"
)

type Middleware = func(http.Handler) http.HandlerFunc

func ApplyMiddlewares(middlewares []Middleware, handlers map[string]http.HandlerFunc) {
	for i := range middlewares {
		for key, value := range handlers {
			handlers[key] = middlewares[i](http.HandlerFunc(value))
		}
	}
}

func (app *App) MiddlewareAuthenticationRequired() Middleware {
	return func(next http.Handler) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := app.Store.Get(r, "user-session")
			if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (app *App) MiddlewareDeveloperOnly() Middleware {
	return func(next http.Handler) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := app.Store.Get(r, "user-session")
			user, _ := session.Values["user"].(database.User)

			if !user.IsDeveloper() {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (app *App) MiddlewareCORS() Middleware {
	return func(next http.Handler) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			for _, allowed := range app.RunningConfiguration.AllowedOrigins {
				if origin == allowed {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Credentials", "true")
					break
				}
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// Preflight
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
