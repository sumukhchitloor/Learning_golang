package main

import (
	"errors"
	"net/http"
	"net/url"
)

var AuthError = errors.New("authentication error")

func Authorize(r *http.Request) error {
	username := r.FormValue("username")
	user, ok := users[username]
	if !ok {
		return AuthError
	}

	// Check session token
	sessionToken, err := r.Cookie("session_token")
	if err != nil || sessionToken.Value != user.SessionToken {
		return AuthError
	}

	// Get CSRF token from header and decode it
	csrfToken := r.Header.Get("X-CSRF-Token")

	// URL decode to handle %3D → = conversion
	decodedCSRF, err := url.QueryUnescape(csrfToken)
	if err != nil {
		// If decoding fails, use original value
		decodedCSRF = csrfToken
	}

	if decodedCSRF != user.CSRFToken || decodedCSRF == "" {
		return AuthError
	}

	return nil
}
