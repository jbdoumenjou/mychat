package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/jbdoumenjou/mychat/internal/log"
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("answer request", slog.String("path", r.URL.Path))

		if _, err := fmt.Fprintln(w, "Hello, World"); err != nil {
			slog.Error("failed to write response", slog.String("error", err.Error()))
		}
	})

	if err := http.ListenAndServe(":8080", nil); err != nil { //nolint:gosec // Reason: This is a simple example.
		slog.Error("failed to start server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
