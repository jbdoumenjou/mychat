package log

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCallerHandler(t *testing.T) {
	var buf bytes.Buffer

	handlerOptions := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	baseHandler := slog.NewJSONHandler(&buf, handlerOptions)
	handler := NewCallerHandler(baseHandler)
	logger := slog.New(handler)

	logger.Info("Test info message")

	output := buf.String()

	if strings.Contains(output, `"caller":`) {
		t.Errorf("did not expect caller information in info log, got: %s", output)
	}

	buf.Reset()
	logger.Debug("Test debug message")

	output = buf.String()
	if !strings.Contains(output, `"caller":`) {
		t.Errorf("expected caller information in debug log, got: %s", output)
	}
}

func TestInitLogger(t *testing.T) {
	var buf bytes.Buffer
	err := InitLogger(&buf, "DEBUG")
	require.NoError(t, err)

	slog.Info("Test info message")

	output := buf.String()

	if strings.Contains(output, `"caller":`) {
		t.Errorf("did not expect caller information in info log, got: %s", output)
	}

	buf.Reset()
	slog.Debug("Test debug message")

	output = buf.String()

	if !strings.Contains(output, `"caller":`) {
		t.Errorf("expected caller information in debug log, got: %s", output)
	}
}
