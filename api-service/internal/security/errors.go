package security

import "errors"

var (
	ErrHashingPassword = errors.New("error hashing your password")
	ErrInvalidPassword = errors.New("invalid password")
	ErrExpiredToken    = errors.New("token has expired")
	ErrInvalidToken    = errors.New("token is invalid")
	ErrUserIDMismatch  = errors.New("context payload doesn't match with header payload")
	ErrPayloadNotFound = errors.New(AuthorizationPayloadKey + " was not found in the request context")
)
