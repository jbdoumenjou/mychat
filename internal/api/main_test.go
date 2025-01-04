package api

import (
	"net/http"
	"testing"

	"github.com/jbdoumenjou/mychat/internal/repo"
)

var (
	testRouter   http.Handler
	testUserRepo *repo.UserRepository
)

func TestMain(m *testing.M) {
	testUserRepo = repo.NewUserRepository()
	messageRepository := repo.NewMessageRepository()
	chatRepository := repo.NewChatRepository()

	userHandler := NewUserHandler(testUserRepo)
	messageHandler := NewMessageHandler(testUserRepo, messageRepository, chatRepository)

	// Create the testRouter
	testRouter = NewRouter(userHandler, messageHandler)

	m.Run()
}
