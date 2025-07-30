package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// ErrInvalidToken is returned when token verification fails
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// Payload contains the payload data of a token including user information and timing
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload creates a new token payload with specified username, role and duration.
// It generates a unique token ID and sets the issued/expired timestamps.
func NewPayload(username string, role string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		Role:      role,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

// Valid checks if the token payload is still valid by verifying the expiration time
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

// JWT-compatible methods for payload to implement jwt.Claims interface

// GetExpirationTime returns the expiration time for JWT compatibility
func (payload *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: payload.ExpiredAt,
	}, nil
}

// GetIssuedAt returns the issued at time for JWT compatibility
func (payload *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: payload.IssuedAt,
	}, nil
}

// GetNotBefore returns the not before time for JWT compatibility
func (payload *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: payload.IssuedAt,
	}, nil
}

// GetIssuer returns empty issuer for JWT compatibility
func (payload *Payload) GetIssuer() (string, error) {
	return "", nil
}

// GetSubject returns empty subject for JWT compatibility
func (payload *Payload) GetSubject() (string, error) {
	return "", nil
}

// GetAudience returns empty audience for JWT compatibility
func (payload *Payload) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{}, nil
}
