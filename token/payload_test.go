package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestNewPayload(t *testing.T) {
	username := "testuser"
	role := "depositor"
	duration := time.Minute

	payload, err := NewPayload(username, role, duration)
	require.NoError(t, err)
	require.NotNil(t, payload)

	require.NotEmpty(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.Equal(t, role, payload.Role)
	require.WithinDuration(t, time.Now(), payload.IssuedAt, time.Second)
	require.WithinDuration(t, time.Now().Add(duration), payload.ExpiredAt, time.Second)
}

func TestPayloadValid(t *testing.T) {
	payload, err := NewPayload("testuser", "depositor", time.Minute)
	require.NoError(t, err)
	require.NotNil(t, payload)

	err = payload.Valid()
	require.NoError(t, err)
}

func TestPayloadExpired(t *testing.T) {
	payload, err := NewPayload("testuser", "depositor", -time.Minute)
	require.NoError(t, err)
	require.NotNil(t, payload)

	err = payload.Valid()
	require.Error(t, err)
	require.Equal(t, ErrExpiredToken, err)
}

func TestPayloadJWTMethods(t *testing.T) {
	payload, err := NewPayload("testuser", "depositor", time.Minute)
	require.NoError(t, err)
	require.NotNil(t, payload)

	// Test GetExpirationTime
	expTime, err := payload.GetExpirationTime()
	require.NoError(t, err)
	require.Equal(t, payload.ExpiredAt, expTime.Time)

	// Test GetIssuedAt
	issuedAt, err := payload.GetIssuedAt()
	require.NoError(t, err)
	require.Equal(t, payload.IssuedAt, issuedAt.Time)

	// Test GetNotBefore
	notBefore, err := payload.GetNotBefore()
	require.NoError(t, err)
	require.Equal(t, payload.IssuedAt, notBefore.Time)

	// Test GetIssuer
	issuer, err := payload.GetIssuer()
	require.NoError(t, err)
	require.Empty(t, issuer)

	// Test GetSubject
	subject, err := payload.GetSubject()
	require.NoError(t, err)
	require.Empty(t, subject)

	// Test GetAudience
	audience, err := payload.GetAudience()
	require.NoError(t, err)
	require.Equal(t, jwt.ClaimStrings{}, audience)
}