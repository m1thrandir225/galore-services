package token

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

// Description:
// Paseto implementation of the Maker interface
//
// Parameters:
// encryptor: *paseto.V2
// symmetricKey: []byte
type PasetoMaker struct {
	encryptor    *paseto.V2
	symmetricKey []byte
}

// Description:
// Return a new PasetoMaker object
//
// Parameters:
// symmetricKey: string
//
// Return:
// Maker
// error
func NewPasetoMaker(symmetricKey string) (Maker, error) {
	log.Println(len(symmetricKey))
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", chacha20poly1305.KeySize)
	}

	pasetoEncryptor := paseto.NewV2()

	maker := &PasetoMaker{
		encryptor:    pasetoEncryptor,
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

// Description:
// Create a paseto universal token
//
// Parameters:
// userID: uuid.UUID
// duration: time.Duration
//
// Return:
// string
// *Payload
// error
func (maker *PasetoMaker) CreateToken(userId uuid.UUID, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userId, duration)
	if err != nil {
		return "", nil, err
	}

	token, err := maker.encryptor.Encrypt(maker.symmetricKey, payload, nil)

	return token, payload, err
}

// Description:
// Verify a generated paseto token and return the payload from the token
//
// Parameters:
// token: string
//
// Return:
// *Payload,
// error
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.encryptor.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}
	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
