// Package api provides a way to write a HTTP Rest API
package api

import "net/http"

// NewRouter is the router for the API.
func NewRouter(users *UserHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /register", users.RegisterUser)

	return mux
}