package weather

import "errors"

var ErrInvalidCity = errors.New("city is invalid")

var ErrInvalidLimit = errors.New("limit is invalid")
