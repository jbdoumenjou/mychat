package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/jbdoumenjou/mychat/internal/repo"
)

// ChatHandler is the handler for user registration.
type ChatHandler struct {
	chatRepo ChatRepo

	logger *slog.Logger
}

// ChatRepo defines the chat repository.
type ChatRepo interface {
	GetUserChats(user string) ([]repo.Chat, error)
}

// NewChatHandler creates a new ChatHandler.
func NewChatHandler(chatRepo ChatRepo) *ChatHandler {
	logger := slog.With(slog.String("handler", "chat"))
	logger.Info("created handler")

	return &ChatHandler{
		chatRepo: chatRepo,
		logger:   logger,
	}
}

// ChatResponse represents a chat between 2 users.
// This is the response format for the API.
// This avoids to expose the internal Chat struct.
type ChatResponse struct {
	ID           string    `json:"id"`
	Participants []string  `json:"participants"`
	CreatedAt    time.Time `json:"createdAt"`
}

// ListChats list all chats for a user.
func (h *ChatHandler) ListChats(w http.ResponseWriter, r *http.Request) {
	h.logger.DebugContext(r.Context(), "handler list chats for a user", slog.String("path", r.URL.Path))

	// TODO: the phone number could be a sensible data, this information should not be expose in the URL.
	// By example, we can use a token to authenticate the user and get the phone number from the token.
	phoneNumber := r.URL.Query().Get("phoneNumber")

	chats, err := h.chatRepo.GetUserChats(phoneNumber)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "failed to get user chats", slog.String("error", err.Error()))
		http.Error(w, "failed to get user chats", http.StatusInternalServerError)

		return
	}

	result := make([]ChatResponse, 0, len(chats))
	for _, chat := range chats {
		result = append(result, ChatResponse{
			ID:           chat.ID,
			Participants: chat.Participants,
			CreatedAt:    chat.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(result); err != nil {
		h.logger.ErrorContext(r.Context(), "failed to write response", slog.String("error", err.Error()))

		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}