package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/jbdoumenjou/mychat/internal/repo"
)

// UserHandler is the handler for user registration.
type UserHandler struct {
	userRepo UserRepo

	logger *slog.Logger
}

// UserRepo defines the user repository.
type UserRepo interface {
	AddUser(phoneNumber string) error
	IsRegistered(phoneNumber string) bool
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(userRepo UserRepo) *UserHandler {
	logger := slog.With(slog.String("handler", "user"))
	logger.Info("created handler")

	return &UserHandler{
		userRepo: userRepo,
		logger:   logger,
	}
}

// User represents a registered user.
type User struct {
	PhoneNumber string `json:"phoneNumber"`
}

// RegisterUser create a new user.
func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	h.logger.DebugContext(r.Context(), "handler register user", slog.String("path", r.URL.Path))

	var user User

	// Decode the JSON body
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.logger.ErrorContext(r.Context(), "Invalid input")
		// TODO: improve error management.
		http.Error(w, "Invalid input", http.StatusBadRequest)

		return
	}

	if user.PhoneNumber == "" {
		h.logger.ErrorContext(r.Context(), "phoneNumber is required")
		http.Error(w, "phoneNumber is required", http.StatusBadRequest)

		return
	}

	// Check if the phone number is already registered
	if h.userRepo.IsRegistered(user.PhoneNumber) {
		h.logger.ErrorContext(r.Context(),
			"Phone number already registered",
			slog.String("phoneNumber", user.PhoneNumber),
		)
		http.Error(w, "Phone number already registered", http.StatusConflict)

		return
	}

	// Register the phone number
	if err := h.userRepo.AddUser(user.PhoneNumber); err != nil {
		h.logger.ErrorContext(r.Context(),
			"Failed to register user",
			slog.String("error", err.Error()),
		)

		if errors.Is(err, repo.ErrPhoneNumberAlreadyRegistered) {
			http.Error(w, "Phone number already registered", http.StatusConflict)

			return
		}

		http.Error(w, "Failed to register user", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	h.logger.DebugContext(r.Context(),
		"User registered successfully",
		slog.Any("user", user),
	)
}
