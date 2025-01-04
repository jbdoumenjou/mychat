package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a test request
	phoneNumber := generateRandomPhoneNumber()
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"/register",
		strings.NewReader(`{"phoneNumber": "`+phoneNumber+`"}`),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Record the HTTP response
	rr := httptest.NewRecorder()
	testRouter.ServeHTTP(rr, req)

	respContentType := rr.Header().Get("Content-Type")
	assert.Equal(t, "application/json", respContentType)
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Try to register the same user again
	req2, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"/register",
		strings.NewReader(`{"phoneNumber": "`+phoneNumber+`"}`),
	)
	require.NoError(t, err)
	req2.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	testRouter.ServeHTTP(rr, req2)

	respContentType2 := rr.Header().Get("Content-Type")
	assert.Equal(t, "text/plain; charset=utf-8", respContentType2)
	assert.Equal(t, http.StatusConflict, rr.Code)
}

func TestRegisterUser_Errors(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tests := []struct {
		name            string
		payload         string
		expectedErrCode int
	}{
		{
			name:            "Invalid phone number",
			payload:         `{"phoneNumber": ""}`,
			expectedErrCode: http.StatusBadRequest,
		},
		{
			name:            "invalid input",
			payload:         `{"phoneNumber": 123}`,
			expectedErrCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequestWithContext(
				ctx,
				http.MethodPost,
				"/register",
				strings.NewReader(test.payload),
			)
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			testRouter.ServeHTTP(rr, req)

			respContentType := rr.Header().Get("Content-Type")
			assert.Equal(t, "text/plain; charset=utf-8", respContentType)
			assert.Equal(t, test.expectedErrCode, rr.Code)
		})
	}
}
