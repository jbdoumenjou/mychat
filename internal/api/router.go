// Package api provides a way to write a HTTP Rest API
package api

import "net/http"

// NewRouter is the router for the API.
func NewRouter(users *UserHandler, messages *MessageHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /register", users.RegisterUser)
	mux.HandleFunc("POST /messages", messages.SendMessage)

	return mux
}
