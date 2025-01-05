// Package api provides a way to write a HTTP Rest API
package api

import "net/http"

// NewRouter is the router for the API.
func NewRouter(users *UserHandler, messages *MessageHandler, chats *ChatHandler) http.Handler {
	mux := http.NewServeMux()

	// user registration with phone number.
	mux.HandleFunc("POST /register", users.RegisterUser)
	// send a message to another user, it will be associated to a chat.
	mux.HandleFunc("POST /messages", messages.SendMessage)
	// list all chats for a user based on the phone number.
	mux.HandleFunc("GET /chats", chats.ListChats)

	return mux
}
