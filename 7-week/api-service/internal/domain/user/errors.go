package user

import "errors"

var (
	ErrNotFound           = errors.New("пользователь не найден")
	ErrEmailAlreadyExists = errors.New("пользователь с таким email уже существует")
	ErrInvalidID          = errors.New("некорректный идентификатор пользователя")
	ErrInvalidFirstName   = errors.New("некорректное имя пользователя")
	ErrInvalidLastName    = errors.New("некорректная фамилия пользователя")
	ErrInvalidEmail       = errors.New("некорректный email пользователя")
	ErrInvalidPassword    = errors.New("некорректный пароль пользователя")
	ErrInvalidRole        = errors.New("некорректная роль пользователя")
)

var fieldErrors = map[string]error{
	"FirstName": ErrInvalidFirstName,
	"LastName":  ErrInvalidLastName,
	"Email":     ErrInvalidEmail,
	"Role":      ErrInvalidRole,
}
