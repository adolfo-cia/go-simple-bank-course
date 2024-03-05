package token

import (
	"errors"
	"time"
)

// Different types of error returned by the VerifyToken function
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// Maker is an interface for managing tokens
type Maker interface {
	// CreateToken create a new token for a specific username and duration
	CreateToken(username string, duration time.Duration) (string, error)
	// VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}
