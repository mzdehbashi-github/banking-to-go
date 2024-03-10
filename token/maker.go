package token

import "time"

// TokenMaker is an interface for managing tokens
type TokenMaker interface {
	// CreateToken creates (and signs) a new token for the given username and duration
	CreateToken(string, time.Duration) (string, error)

	// VerifyToken Validates if the token is valid or not
	VerifyToken(string) (*Payload, error)
}
