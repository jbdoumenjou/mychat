// Package repo manage storage and operations.
package repo

import (
	"log/slog"
	"sync"
)

// MessageRepository manages user storage and operations
// In-memory store for simplicity
// TODO: use a database instead.
type MessageRepository struct {
	mu       sync.Mutex
	messages map[string][]string // map[chatID][message1, message2, message3]

	logger *slog.Logger
}

// NewMessageRepository initializes a new MessageRepository.
func NewMessageRepository() *MessageRepository {
	logger := slog.With(slog.String("repo", "message"))
	logger.Info("created repository")

	return &MessageRepository{
		messages: make(map[string][]string),
		logger:   logger,
	}
}

// AddMessage adds a new message to the repository.
func (repo *MessageRepository) AddMessage(chatID, content string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.messages[chatID] = append(repo.messages[chatID], content)

	return nil
}