package history

import "errors"

var (
	ErrInvalidID     = errors.New("некорректный идентификатор записи истории")
	ErrInvalidUserID = errors.New("некорректный идентификатор пользователя")
	ErrInvalidCity   = errors.New("некорректный город")
)
