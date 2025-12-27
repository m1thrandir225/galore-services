package token

import (
	"time"

	"github.com/google/uuid"
)

// Maker represents a generic interface that all token providers need to implement
type Maker interface {

	//CreateToken creates a token and a payload
	CreateToken(userId uuid.UUID, duration time.Duration) (string, *Payload, error)

	//VerifyToken verifies if the given token is valid
	VerifyToken(token string) (*Payload, error)
}
