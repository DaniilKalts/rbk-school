package history

import "errors"

var (
	ErrInvalidID     = errors.New("history id is invalid")
	ErrInvalidUserID = errors.New("user id is invalid")
	ErrInvalidCity   = errors.New("city is invalid")
)
