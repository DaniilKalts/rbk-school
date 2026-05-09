package history

import "errors"

var (
	ErrInvalidUserID = errors.New("некорректный идентификатор пользователя")
	ErrInvalidCity   = errors.New("некорректный город")
)
