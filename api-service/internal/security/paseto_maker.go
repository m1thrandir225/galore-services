package security

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

// PASETOMaker implements the Maker interface for PASETO tokens
type PASETOMaker struct {
	encryptor    *paseto.V2
	symmetricKey []byte
}

// NewPASETOMaker returns a PASETOMaker instance
func NewPASETOMaker(symmetricKey string) (*PASETOMaker, error) {
	log.Println(len(symmetricKey))
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", chacha20poly1305.KeySize)
	}

	pasetoEncryptor := paseto.NewV2()

	maker := &PASETOMaker{
		encryptor:    pasetoEncryptor,
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

func (maker *PASETOMaker) CreateToken(userId uuid.UUID, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userId, duration)
	if err != nil {
		return "", nil, err
	}

	token, err := maker.encryptor.Encrypt(maker.symmetricKey, payload, nil)

	return token, payload, err
}

func (maker *PASETOMaker) VerifyToken(token string) (*Payload, error) {
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
