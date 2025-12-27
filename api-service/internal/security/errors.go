package security

import "errors"

var (
	ErrHashingPassword = errors.New("error hashing your password")
	ErrInvalidPassword = errors.New("invalid password")
)
