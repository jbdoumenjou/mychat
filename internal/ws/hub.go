// Package ws contains tools to use webservices.
package ws

import (
	"log/slog"
	"sync"
)

// Message represents a message sent between clients.
type Message struct {
	From    string
	To      string
	Content string
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[string]*Client
	mu      sync.RWMutex
}

// NewHub creates a new Hub.
func NewHub() *Hub {
	slog.Info("created hub")

	return &Hub{
		clients: make(map[string]*Client),
	}
}

// Register adds a client to the hub.
func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client.UserID]; ok {
		slog.Warn("client already registered", slog.String("client", client.UserID))

		return
	}

	h.clients[client.UserID] = client
	slog.Info("client registered", slog.String("client", client.UserID))
}

// Unregister removes a client to the hub.
func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client.UserID]; !ok {
		slog.Warn("client not registered", slog.String("client", client.UserID))

		return
	}

	delete(h.clients, client.UserID)
	slog.Info("client unregistered", slog.String("client", client.UserID))
}

// Broadcast sends a message to all targeted clients.
func (h *Hub) Broadcast(msg *Message) {
	slog.Info("broadcasting message", slog.Any("msg", msg))

	h.mu.RLock()
	defer h.mu.RUnlock()

	client, ok := h.clients[msg.To]
	if !ok {
		slog.Warn("client not found", slog.String("client", msg.To))

		return
	}

	if err := client.conn.WriteJSON(msg); err != nil {
		slog.Error("failed to send message", slog.String("err", err.Error()))

		return
	}
}
