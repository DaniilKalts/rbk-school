package weather

import "errors"

var ErrInvalidCity = errors.New("city is invalid")

var ErrInvalidLimit = errors.New("limit is invalid")

var ErrInvalidOffset = errors.New("offset is invalid")
