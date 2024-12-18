package token

import (
	"github.com/google/uuid"
	"time"
)

type Maker interface {
	CreateToken(userId uuid.UUID, duration time.Duration) (string, interface{}, error)
	VerifyToken(token string) (interface{}, error)
}
