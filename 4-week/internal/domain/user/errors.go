package user

import "errors"

var (
	ErrNotFound           = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("user email already exists")
	ErrInvalidID          = errors.New("user id is invalid")
	ErrInvalidFirstName   = errors.New("user first name is invalid")
	ErrInvalidLastName    = errors.New("user last name is invalid")
	ErrInvalidEmail       = errors.New("user email is invalid")
	ErrInvalidPassword    = errors.New("user password is invalid")
	ErrInvalidRole        = errors.New("user role is invalid")
)

var fieldErrors = map[string]error{
	"FirstName": ErrInvalidFirstName,
	"LastName":  ErrInvalidLastName,
	"Email":     ErrInvalidEmail,
	"Role":      ErrInvalidRole,
}
