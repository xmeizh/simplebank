package token

import (
	"errors"
	"time"
)

var ErrInvalidToken = errors.New("token is invalid")
var ErrTokenExpired = errors.New("token is expired")

// Maker is an interface for managing tokens
type Maker interface {
	CreateToken(username string, duration time.Duration) (string, *Payload, error)

	// VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}
