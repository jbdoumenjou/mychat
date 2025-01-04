// Package log provides a logger ability like caller information to DEBUG logs.
package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"strconv"
)

type callerHandler struct {
	wrapped slog.Handler
}

func (h *callerHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.wrapped.Enabled(ctx, level)
}

//nolint:gocritic // Reason: Following the wrapped implementation which passes record by value.
func (h *callerHandler) Handle(ctx context.Context, record slog.Record) error {
	// Only add caller information for DEBUG logs
	if record.Level == slog.LevelDebug {
		// Adjust the skip level if necessary
		if _, file, line, ok := runtime.Caller(3); ok {
			record.AddAttrs(slog.String("caller", file+":"+strconv.Itoa(line)))
		}
	}

	return fmt.Errorf("failed to handle log record: %w", h.wrapped.Handle(ctx, record))
}

func (h *callerHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &callerHandler{
		wrapped: h.wrapped.WithAttrs(attrs),
	}
}

func (h *callerHandler) WithGroup(name string) slog.Handler {
	return &callerHandler{
		wrapped: h.wrapped.WithGroup(name),
	}
}

// NewCallerHandler wraps the given handler and adds caller information to DEBUG logs.
func NewCallerHandler(wrapped slog.Handler) slog.Handler {
	return &callerHandler{wrapped: wrapped}
}

// InitLogger initializes the logger with the given log level.
func InitLogger(w io.Writer, level string) error {
	var logLevel slog.Level
	if err := logLevel.UnmarshalText([]byte(level)); err != nil {
		slog.Error("failed to parse log level", slog.String("error", err.Error()))

		return fmt.Errorf("failed to parse log level: %w", err)
	}

	handlerOptions := &slog.HandlerOptions{
		Level: logLevel,
	}

	baseHandler := slog.NewJSONHandler(w, handlerOptions)
	logger := slog.New(NewCallerHandler(baseHandler))

	slog.SetDefault(logger)

	// Log here to use initialized logger.
	slog.Info("logger initialized",
		slog.String("level", logLevel.String()),
		slog.String("type", "json"),
	)

	return nil
}
