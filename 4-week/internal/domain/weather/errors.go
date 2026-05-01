package weather

import "errors"

var ErrInvalidCity = errors.New("некорректный город")

var ErrInvalidLimit = errors.New("некорректный limit")

var ErrInvalidOffset = errors.New("некорректный offset")
