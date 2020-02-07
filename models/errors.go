package models

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserPhoneExists    = errors.New("user with this phone number already exists")
	ErrInvalidCredentials = errors.New("invalid login credentials")
	ErrInvalidRegCode     = errors.New("registration code is invalid or used")
)
