package city

import "errors"

var (
	ErrNotFound      = errors.New("city not found")
	ErrAlreadyExists = errors.New("city already exists")
	ErrInvalidID     = errors.New("city id is invalid")
	ErrInvalidUserID = errors.New("user id is invalid")
	ErrInvalidName   = errors.New("city name is invalid")
)
