package ws

import (
	"github.com/gorilla/websocket"
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	UserID string

	// The websocket connection.
	conn *websocket.Conn
}

// NewClient creates a new client.
func NewClient(conn *websocket.Conn, userID string) *Client {
	return &Client{
		UserID: userID,
		conn:   conn,
	}
}
