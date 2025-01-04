package api

import (
	"net/http"
	"testing"

	"github.com/jbdoumenjou/mychat/internal/repo"
)

var testRouter http.Handler

func TestMain(m *testing.M) {
	userRepository := repo.NewUserRepository()
	userHandler := NewUserHandler(userRepository)

	// Create the testRouter
	testRouter = NewRouter(userHandler)

	m.Run()
}
