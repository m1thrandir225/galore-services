package services

import "errors"

var (
	ErrorCreatingHOTPCounter = errors.New("error creating HOTP counter")
	ErrorGettingHOTPCounter  = errors.New("error getting HOTP counter")
)
