package city

import "errors"

var (
	ErrNotFound      = errors.New("город не найден")
	ErrAlreadyExists = errors.New("город уже существует")
	ErrInvalidID     = errors.New("некорректный идентификатор города")
	ErrInvalidUserID = errors.New("некорректный идентификатор пользователя")
	ErrInvalidName   = errors.New("некорректное название города")
)
