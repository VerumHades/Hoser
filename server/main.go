package main

import (
	"fmt"
	"net/http"

	"server/internal/database"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte("super-secret-key"))
var database_interactor database.Interactor = &database.DummyInteractor{}

func sendIvalidCredentialsError(w http.ResponseWriter) {
	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	user, user_error := database_interactor.GetUserByName(username)

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

	fmt.Fprintln(w, "Logged in")
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	session, _ := store.Get(r, "user-session")

	session.Values["authenticated"] = false
	session.Options.MaxAge = -1
	session.Save(r, w)

	fmt.Fprintln(w, "Logged out")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	session, _ := store.Get(r, "user-session")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	user := session.Values["user"]
	fmt.Fprintf(w, "Welcome, %s!", user)
}

var clientBuildDirectory = "../client/dist"

func main() {
	// API endpoints
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/home", homeHandler)

	fs := http.FileServer(http.Dir(clientBuildDirectory))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := http.Dir(clientBuildDirectory).Open(r.URL.Path); err != nil {
			http.ServeFile(w, r, clientBuildDirectory+"/index.html")
			return
		}
		fs.ServeHTTP(w, r)
	})

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe("localhost:8080", nil)
}
