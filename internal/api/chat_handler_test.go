package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChatHandler_ListChats(t *testing.T) {
	users := make([]string, 0, 10)

	for range 10 {
		phoneNumber := generateRandomPhoneNumber()
		users = append(users, phoneNumber)
		err := testUserRepo.AddUser(phoneNumber)
		require.NoError(t, err)
	}

	var chats []string

	for i := 1; i < len(users); i++ {
		chat, err := testChatRepo.GetOrCreateChat(users[0], users[i])
		require.NoError(t, err)

		chats = append(chats, chat)
	}

	// Create a test request
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/chats?phoneNumber="+url.QueryEscape(users[0]),
		http.NoBody,
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Record the HTTP response
	rr := httptest.NewRecorder()
	testRouter.ServeHTTP(rr, req)

	respContentType := rr.Header().Get("Content-Type")
	assert.Equal(t, "application/json", respContentType)
	assert.Equal(t, http.StatusOK, rr.Code)

	var result []ChatResponse
	err = json.NewDecoder(rr.Body).Decode(&result)
	require.NoError(t, err)
	require.NotEmpty(t, result)
	// TODO: improve the test by checking the content of the result.
	require.Len(t, chats, 9)
}

// TODO: add errors tests for ListChats

func TestChatHandler_ListChatMessages(t *testing.T) {
	users := make([]string, 0, 10)

	for range 2 {
		phoneNumber := generateRandomPhoneNumber()
		users = append(users, phoneNumber)
		err := testUserRepo.AddUser(phoneNumber)
		require.NoError(t, err)
	}

	chatID, err := testChatRepo.GetOrCreateChat(users[0], users[1])
	require.NoError(t, err)

	msgs := make([]string, 0, 3)
	for i := range 3 {
		msgs = append(msgs, fmt.Sprintf("Content msg {%d}\n", i))

		require.NoError(t, testMessageRepo.AddMessage(chatID, msgs[i]))
	}

	// Create a test request
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/chats/"+chatID+"/messages",
		http.NoBody,
	)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	// Record the HTTP response
	rr := httptest.NewRecorder()
	testRouter.ServeHTTP(rr, req)

	respContentType := rr.Header().Get("Content-Type")
	assert.Equal(t, "application/json", respContentType)
	assert.Equal(t, http.StatusOK, rr.Code)

	var result []string
	err = json.NewDecoder(rr.Body).Decode(&result)
	require.NoError(t, err)
	require.NotEmpty(t, result)
	require.Equal(t, msgs, result)
}

// TODO: add errors tests for ListChatMessages
