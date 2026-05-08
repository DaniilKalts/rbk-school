package weather

import "errors"

var ErrInvalidCity = errors.New("некорректный город")

var ErrInvalidLatitude = errors.New("некорректная широта")

var ErrInvalidLongitude = errors.New("некорректная долгота")

var ErrInvalidLimit = errors.New("некорректный limit")

var ErrInvalidOffset = errors.New("некорректный offset")
