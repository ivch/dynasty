package models

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserPhoneExists    = errors.New("user with this phone number already exists")
	ErrInvalidCredentials = errors.New("invalid login credentials")
	ErrInvalidRegCode     = errors.New("registration code is invalid or used")

	ErrParsingToken       = errors.New("failed to parse token")
	ErrParsingTokenClaims = errors.New("failed to parse token claims")
	ErrTokenIsInvalid     = errors.New("token is invalid")
	ErrTokenExpired       = errors.New("token expired")

	ErrSessionNotFound = errors.New("session not found")
)
