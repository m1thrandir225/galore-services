// Description:
// Security package for Galore. Mainly for OTP generation & validation
package security

import (
	"crypto/rand"
	"encoding/base32"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/hotp"
)

// Description:
// Generates a hash based OTP code
//
// Parameters:
// secret: string
// counter: uin64
//
// Return:
// string
// error
func GenerateHOTP(secret string, counter uint64) (string, error) {
	opts := hotp.ValidateOpts{
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA256,
	}

	code, err := hotp.GenerateCodeCustom(secret, counter, opts)
	if err != nil {
		return "", err
	}
	return code, nil
}

// Description:
// A method that validates a given hash based OTP code
//
// Parameters:
// secret: string
// providedOTP: string
// counter: uint64
//
// Return:
// bool
// error
func ValidateHOTP(secret string, providedOTP string, counter uint64) (bool, error) {
	opts := hotp.ValidateOpts{
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA256,
	}

	valid, err := hotp.ValidateCustom(providedOTP, counter, secret, opts)
	if err != nil {
		return false, err
	}

	return valid, nil
}

// Description:
// Generates a unique OTP secret for a user
//
// Return:
// string
// error
func GenerateOTPSecret() (string, error) {
	secret := make([]byte, 20)

	_, err := rand.Read(secret)
	if err != nil {
		return "", err
	}

	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(secret), nil
}
