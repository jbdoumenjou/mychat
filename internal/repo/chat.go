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
	ID        string
	CreatedAt time.Time
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
func (repo *ChatRepository) GetOrCreateChat(sender, receiver string) (string, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.chats[sender]; !exists {
		repo.chats[sender] = make(map[string]*Chat)
	}

	if _, exists := repo.chats[receiver]; !exists {
		repo.chats[receiver] = make(map[string]*Chat)
	}

	chat, exists := repo.chats[sender][receiver]
	if !exists || chat == nil {
		chat = &Chat{
			ID:        uuid.NewString(),
			CreatedAt: time.Now().UTC().Truncate(time.Millisecond),
		}

		// add chat to sender and receiver to retrieve all chat for a user
		repo.chats[sender][receiver] = chat
		repo.chats[receiver][sender] = chat

		// Add to `chatsByUser`
		repo.chatsByUser[sender] = append(repo.chatsByUser[sender], chat)
		repo.chatsByUser[receiver] = append(repo.chatsByUser[receiver], chat)
	}

	return chat.ID, nil
}

// GetChatsByUser gets all chats for a user.
func (repo *ChatRepository) GetChatsByUser(user string) ([]*Chat, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	return repo.chatsByUser[user], nil
}
