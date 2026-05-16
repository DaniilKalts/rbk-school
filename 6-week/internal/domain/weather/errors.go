package weather

import "errors"

var (
	ErrInvalidCity        = errors.New("некорректный город")
	ErrInvalidTemperature = errors.New("некорректная температура")
	ErrInvalidLimit       = errors.New("некорректный limit")
	ErrInvalidOffset      = errors.New("некорректный offset")
)
