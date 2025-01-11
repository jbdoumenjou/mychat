package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jbdoumenjou/mychat/internal/ws"
)

// MessageHandler is the handler for user registration.
type MessageHandler struct {
	chatRepo    MessageChatRepo
	messageRepo MessageRepo
	userRepo    MessageUserRepo
	hub         *ws.Hub
	upgrader    websocket.Upgrader

	logger *slog.Logger
}

// MessageUserRepo defines the user repository.
type MessageUserRepo interface {
	IsRegistered(phoneNumber string) bool
}

// MessageRepo defines the user repository.
type MessageRepo interface {
	AddMessage(chatID, content string) error
}

// MessageChatRepo defines the chat repository.
type MessageChatRepo interface {
	GetOrCreateChat(sender, receiver string) (string, error)
}

// NewMessageHandler creates a new MessageHandler.
func NewMessageHandler(userRepo MessageUserRepo, messageRepo MessageRepo, chatRepo MessageChatRepo) *MessageHandler {
	logger := slog.With(slog.String("handler", "message"))
	logger.Info("created handler")

	return &MessageHandler{
		chatRepo:    chatRepo,
		messageRepo: messageRepo,
		userRepo:    userRepo,
		// TODO: Check if the hub should be injected or created here.
		hub: ws.NewHub(),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(_ *http.Request) bool {
				return true // Allow all origins for simplicity; adjust for security.
			},
		},
		logger: logger,
	}
}

// Message represents a message to send from a user to another user.
type Message struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Content  string `json:"content"`
}

// SendMessage create a new message in a chat with two users.
func (h *MessageHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	h.logger.DebugContext(r.Context(), "handler register message", slog.String("path", r.URL.Path))

	var message Message

	// Decode the JSON body
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		h.logger.ErrorContext(r.Context(), "Invalid input")
		// TODO: improve error management.
		http.Error(w, "Invalid input", http.StatusBadRequest)

		return
	}

	if message.Content == "" {
		h.logger.ErrorContext(r.Context(), "message content is required")
		http.Error(w, "message content is required", http.StatusBadRequest)

		return
	}

	// Check if the phone number is already registered.
	if !h.userRepo.IsRegistered(message.Sender) {
		h.logger.ErrorContext(r.Context(),
			"Sender phone number not registered",
			slog.String("phoneNumber", message.Sender),
		)
		http.Error(w, "Sender phone number not registered", http.StatusBadRequest)

		return
	}

	// Check if the phone number is already registered.
	if !h.userRepo.IsRegistered(message.Receiver) {
		h.logger.ErrorContext(r.Context(),
			"Receiver phone number not registered",
			slog.String("phoneNumber", message.Receiver),
		)
		http.Error(w, "Receiver phone number not registered", http.StatusBadRequest)

		return
	}

	// Get or create the chat with the two users.
	chatID, err := h.chatRepo.GetOrCreateChat(message.Sender, message.Receiver)
	if err != nil {
		h.logger.ErrorContext(r.Context(),
			"Failed to create chatID",
			slog.String("error", err.Error()),
		)
		http.Error(w, "Failed to create chatID", http.StatusInternalServerError)

		return
	}

	// Add the message to the chat.
	if err = h.messageRepo.AddMessage(chatID, message.Content); err != nil {
		h.logger.ErrorContext(r.Context(),
			"Failed to register message",
			slog.String("error", err.Error()),
		)
		http.Error(w, "Failed to register message", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	h.logger.DebugContext(r.Context(),
		"Content send successfully",
		slog.Any("message", message),
	)
}

// HandleWS upgrades the HTTP connection to a WebSocket and manages messages.
func (h *MessageHandler) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("failed to upgrade to WebSocket", slog.String("error", err.Error()))

		return
	}

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		slog.Error("missing userID in query params")

		return
	}

	client := ws.NewClient(conn, userID)
	h.hub.Register(client)

	if err = h.listenWS(r.Context(), conn); err != nil {
		slog.Error("failed to listenWS", slog.String("error", err.Error()))
		http.Error(w, "failed to listenWS", http.StatusInternalServerError)

		return
	}
}

// listenWS reads messages from the WebSocket connection.
func (h *MessageHandler) listenWS(ctx context.Context, conn *websocket.Conn) error {
	for {
		select {
		case <-ctx.Done():
			slog.Info("stop listening websocket connection")

			return nil
		default:
			var msg ws.Message
			if err := conn.ReadJSON(&msg); err != nil {
				slog.Error("failed to read message", slog.String("error", err.Error()))

				return fmt.Errorf("failed to read message: %w", err)
			}

			h.hub.Broadcast(&msg)
		}
	}
}
