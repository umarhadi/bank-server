// Package token provides token generation and verification functionality
// for user authentication using both JWT and PASETO token formats.
package token

import "time"

// Maker defines the interface for token generation and verification.
// It supports creating tokens with username, role, and duration,
// and verifying tokens to extract payload information.
type Maker interface {
	// CreateToken creates a new token for a specific username and role with given duration
	CreateToken(username string, role string, duration time.Duration) (string, *Payload, error)

	// VerifyToken checks a token string and returns the payload if valid
	VerifyToken(token string) (*Payload, error)
}
