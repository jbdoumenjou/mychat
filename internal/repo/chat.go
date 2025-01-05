// Package repo manage storage and operations.
package repo

import (
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Chat represents a chat between 2 users.
type Chat struct {
	ID           string
	Participants []string // user IDs
	CreatedAt    time.Time
}

// ChatRepository manages chat storage and operations
// In-memory store for simplicity
// TODO: use a database instead.
type ChatRepository struct {
	mu          sync.RWMutex
	chats       map[string]map[string]*Chat // map[user][user2]Chat
	chatsByUser map[string][]*Chat

	logger *slog.Logger
}

// NewChatRepository initializes a new ChatRepository.
func NewChatRepository() *ChatRepository {
	logger := slog.With(slog.String("repo", "chat"))
	logger.Info("created repository")

	return &ChatRepository{
		chats:       make(map[string]map[string]*Chat),
		chatsByUser: make(map[string][]*Chat),
		logger:      logger,
	}
}

// GetOrCreateChat gets a chat from the repository.
// TODO: refactor to avoid doing 2 things in one function.
// This is a very naive approach to chat management.
func (r *ChatRepository) GetOrCreateChat(sender, receiver string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.chats[sender]; !exists {
		r.chats[sender] = make(map[string]*Chat)
	}

	if _, exists := r.chats[receiver]; !exists {
		r.chats[receiver] = make(map[string]*Chat)
	}

	chat, exists := r.chats[sender][receiver]
	if !exists || chat == nil {
		chat = &Chat{
			ID:           uuid.NewString(),
			Participants: []string{sender, receiver},
			CreatedAt:    time.Now().UTC().Truncate(time.Millisecond),
		}

		// add chat to sender and receiver to retrieve all chat for a user
		r.chats[sender][receiver] = chat
		r.chats[receiver][sender] = chat

		// Add to `chatsByUser`
		r.chatsByUser[sender] = append(r.chatsByUser[sender], chat)
		r.chatsByUser[receiver] = append(r.chatsByUser[receiver], chat)
	}

	return chat.ID, nil
}

// GetUserChats gets all chats for a user.
// TODO:
// * add pagination. (more appropriate for a database)
// * return a DTO instead of the entity.
func (r *ChatRepository) GetUserChats(user string) ([]Chat, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if _, exists := r.chatsByUser[user]; !exists {
		return []Chat{}, nil
	}

	// dereference the chats
	result := make([]Chat, len(r.chatsByUser[user]))
	for i, chat := range r.chatsByUser[user] {
		result[i] = *chat
	}

	r.logger.Debug("get user chats",
		slog.String("user", user),
		slog.Any("chats", result),
	)

	return result, nil
}
