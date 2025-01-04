package repo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserRepository_AddUser(t *testing.T) {
	userRepo := NewUserRepository()

	err := userRepo.AddUser("123")
	require.NoError(t, err)

	err = userRepo.AddUser("123")
	require.ErrorIs(t, err, ErrPhoneNumberAlreadyRegistered)

	err = userRepo.AddUser("1234")
	require.NoError(t, err)
}
