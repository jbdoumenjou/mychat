package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jbdoumenjou/mychat/internal/api"
	"github.com/jbdoumenjou/mychat/internal/log"
	"github.com/jbdoumenjou/mychat/internal/repo"
)

func main() {
	// Get the log level from the environment variable
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "INFO" // Default log level
	}

	// Initialize logger.
	if err := log.InitLogger(os.Stdout, logLevel); err != nil {
		slog.Error("failed to initialize logger", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Repository
	userRepo := repo.NewUserRepository()
	chatRepo := repo.NewChatRepository()
	messageRepo := repo.NewMessageRepository()

	// API handlers
	userHandler := api.NewUserHandler(userRepo)
	messageHandler := api.NewMessageHandler(userRepo, messageRepo, chatRepo)
	router := api.NewRouter(userHandler, messageHandler)

	// Create an HTTP server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		IdleTimeout:  10 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Channel to listen for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Run the server in a goroutine
	go func() {
		slog.Info("Starting server on :8080")

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Error starting server", slog.String("error", err.Error()))
		}
	}()

	// Wait for a signal
	<-stop
	slog.Info("Shutting down server...")

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", slog.String("error", err.Error()))
	} else {
		slog.Info("Server shutdown gracefully")
	}
}
