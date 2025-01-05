package api

import (
	"net/http"
	"testing"

	"github.com/jbdoumenjou/mychat/internal/repo"
)

var (
	testRouter      http.Handler
	testUserRepo    *repo.UserRepository
	testChatRepo    *repo.ChatRepository
	testMessageRepo *repo.MessageRepository
)

func TestMain(m *testing.M) {
	testUserRepo = repo.NewUserRepository()
	testMessageRepo = repo.NewMessageRepository()
	testChatRepo = repo.NewChatRepository()

	userHandler := NewUserHandler(testUserRepo)
	messageHandler := NewMessageHandler(testUserRepo, testMessageRepo, testChatRepo)
	chatHandler := NewChatHandler(testChatRepo)

	// Create the testRouter
	testRouter = NewRouter(userHandler, messageHandler, chatHandler)

	m.Run()
}
