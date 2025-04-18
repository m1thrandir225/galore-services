package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Description:
// Hash a given string using bcrypt
//
// Parameters:
// base: string
//
// Return:
// string
// error
func HashPassowrd(base string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(base), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("there was an error hashing your password: %s", err.Error())
	}

	return string(hashedPassword), err
}

// Description:
// Compare a given unhashed password with it's hash using bcrypt
//
// Parameters:
// hashedPassword: string
// password: string
//
// Return:
// error
func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
