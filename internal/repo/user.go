// Package repo manage storage and operations.
package repo

import (
	"log/slog"
	"sync"
)

// UserRepository manages user storage and operations
// In-memory store for simplicity
// TODO: use a database instead.
type UserRepository struct {
	mu    sync.Mutex
	users map[string]bool

	logger *slog.Logger
}

// NewUserRepository initializes a new UserRepository.
func NewUserRepository() *UserRepository {
	logger := slog.With(slog.String("repo", "user"))
	logger.Info("created repository")

	return &UserRepository{
		users:  make(map[string]bool),
		logger: logger,
	}
}

// AddUser adds a new user to the repository.
func (repo *UserRepository) AddUser(phoneNumber string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.users[phoneNumber] {
		return ErrPhoneNumberAlreadyRegistered
	}

	repo.users[phoneNumber] = true

	return nil
}

// IsRegistered checks if a phone number is already registered.
func (repo *UserRepository) IsRegistered(phoneNumber string) bool {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	return repo.users[phoneNumber]
}
