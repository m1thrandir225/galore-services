package security

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a given string using bcrypt
func HashPassword(base string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(base), bcrypt.DefaultCost)
	if err != nil {
		return "", ErrHashingPassword
	}

	return string(hashedPassword), err
}

// ComparePassword compares a hashed and unhashed string using  bcrypt
func ComparePassword(hashedPassword, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return ErrInvalidPassword
	}
	return nil
}
