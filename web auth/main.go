package main

import (
	"fmt"
	"net/http"
	"time"
)

type Login struct {
	HashedPassword string
	SessionToken   string
	CSRFToken      string
}

// key is the username
var users = map[string]Login{}

func main() {
	http.HandleFunc("/register", register)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/login", login)
	http.HandleFunc("/protected", protected)

	http.ListenAndServe(":8080", nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		er := http.StatusMethodNotAllowed
		http.Error(w, "Invalid request method", er)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	// Validate input
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	if len(password) < 8 {
		http.Error(w, "Password must be at least 8 characters long", http.StatusBadRequest)
		return
	}

	if _, ok := users[username]; ok {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	// Hash password and handle error
	hashedPassword, err := hashPassword(password)
	if err != nil {
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	users[username] = Login{
		HashedPassword: hashedPassword,
	}

	// Send proper response to client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"message": "User registered successfully", "username": "%s"}`, username)
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		er := http.StatusMethodNotAllowed
		http.Error(w, "Invalid request method", er)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	user, ok := users[username]
	if !ok || !checkPasswordHash(password, user.HashedPassword) {
		// If user does not exist or password is incorrect
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		// Return a JSON error response
		er := http.StatusUnauthorized
		http.Error(w, "Invalid username or password", er)
		return
	}

	sessionToken := generateToken(32) // Generate a session token
	csrfToken := generateToken(32)    // Generate a CSRF token

	// set session token and CSRF token in the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(24 * time.Hour), // Set cookie to expire in 24 hours
		HttpOnly: true,                           // Make cookie HTTP-only for security

	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  time.Now().Add(24 * time.Hour), // Set cookie to expire in 24 hours
		HttpOnly: false,                          // Needs to be accessible to the client-side for CSRF protection
	})

	user.SessionToken = sessionToken
	user.CSRFToken = csrfToken
	users[username] = user // Update the user with session and CSRF tokens

	fmt.Fprintf(w, `{"message": "Login successful", "username": "%s"}`, username)

}

func protected(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if err := Authorize(r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	username := r.FormValue("username")
	fmt.Fprintf(w, "CSRF token is valid for user: %s", username)
}

func logout(w http.ResponseWriter, r *http.Request) {
	if err := Authorize(r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	//clear session token and CSRF token
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Set cookie to expire in the past
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Set cookie to expire in the past
		HttpOnly: false,
	})

	//clear token from the database
	username := r.FormValue("username")
	if user, ok := users[username]; ok {
		user.SessionToken = ""
		user.CSRFToken = ""
		users[username] = user // Update the user in the map
	}

	fmt.Fprintf(w, `{"message": "Logout successful", "username": "%s"}`, username)
}
