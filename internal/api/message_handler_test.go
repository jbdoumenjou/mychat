package api

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendMessage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// add a sender
	sender := generateRandomPhoneNumber()
	err := testUserRepo.AddUser(sender)
	require.NoError(t, err)

	// add a receiver
	receiver := generateRandomPhoneNumber()
	err = testUserRepo.AddUser(receiver)
	require.NoError(t, err)

	// Create a test request
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"/messages",
		strings.NewReader(`{"sender": "`+sender+`", "receiver": "`+receiver+`", "content": "Hello World!"}`),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Record the HTTP response
	rr := httptest.NewRecorder()
	testRouter.ServeHTTP(rr, req)

	respContentType := rr.Header().Get("Content-Type")
	assert.Equal(t, "application/json", respContentType)
	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestSendMessage_Errors(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// add a first user
	user1 := generateRandomPhoneNumber()
	err := testUserRepo.AddUser(user1)
	require.NoError(t, err)

	// add a second user
	user2 := generateRandomPhoneNumber()
	err = testUserRepo.AddUser(user2)
	require.NoError(t, err)

	tests := []struct {
		name            string
		sender          string
		receiver        string
		content         string
		expectedCodeErr int
	}{
		{
			name:            "Empty user1",
			sender:          "",
			receiver:        user2,
			content:         "Hello World!",
			expectedCodeErr: http.StatusBadRequest,
		},
		{
			name:            "Empty user2",
			sender:          user1,
			receiver:        "",
			content:         "Hello World!",
			expectedCodeErr: http.StatusBadRequest,
		},
		{
			name:            "Empty content",
			sender:          user1,
			receiver:        user2,
			content:         "",
			expectedCodeErr: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a test request
			req, err := http.NewRequestWithContext(
				ctx,
				http.MethodPost,
				"/messages",
				strings.NewReader(`{"sender": "`+test.sender+`", "receiver": "`+test.receiver+`", "content": "`+test.content+`"}`),
			)
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Record the HTTP response
			rr := httptest.NewRecorder()
			testRouter.ServeHTTP(rr, req)

			respContentType := rr.Header().Get("Content-Type")
			assert.Equal(t, "text/plain; charset=utf-8", respContentType)
			assert.Equal(t, test.expectedCodeErr, rr.Code)
		})
	}
}

// generateRandomPhoneNumber generates a random phone number in E.164 format.
func generateRandomPhoneNumber() string {
	// Generate a random country code (1 to 999)
	countryCode := rand.Intn(999) + 1

	// Generate a random local number (10 digits)
	localNumber := rand.Intn(1000000000) + 1000000000

	// Return the phone number in E.164 format
	return fmt.Sprintf("+%d%d", countryCode, localNumber)
}
