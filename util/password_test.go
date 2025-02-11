package util

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPass(t *testing.T) {
	password := RandomString(8)

	hashedPassword1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	err = CheckPassword(password, hashedPassword1)
	require.NoError(t, err)

	wrongPassword := RandomString(8)
	err = CheckPassword(wrongPassword, hashedPassword1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedPassword2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	require.NotEqual(t, hashedPassword1, hashedPassword2)
}

func TestHashPassword_Error(t *testing.T) {
	originalFunc := generateFromPassword
	defer func() { generateFromPassword = originalFunc }()

	generateFromPassword = func(password []byte, cost int) ([]byte, error) {
		return nil, errors.New("simulated error")
	}

	hashed, err := HashPassword("any-password")
	require.Empty(t, hashed)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to hash password: simulated error")
}
