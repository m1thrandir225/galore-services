package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

var (
	ErrExpiredToken = errors.New("the token has expired")
	ErrInvalidToken = errors.New("the token is invalid")
)

type JWTClaims struct {
	ID                   uuid.UUID `json:"id"`
	UserID               uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims `json:",inline"`
}

func NewJWTClaims(userID uuid.UUID, duration time.Duration) (*JWTClaims, error) {
	tokenId, err := uuid.NewUUID()

	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(duration)

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Subject:   tokenId.String(),
		Issuer:    "galoreapp",
		Audience:  []string{userID.String()},
	}

	payload := &JWTClaims{
		ID:               tokenId,
		UserID:           userID,
		RegisteredClaims: claims,
	}
	return payload, nil
}
