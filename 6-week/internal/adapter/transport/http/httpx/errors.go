package httpx

import (
	"errors"
	"net/http"

	"github.com/DaniilKalts/rbk-school/6-week/internal/adapter/jwt"
	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/city"
	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/user"
	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/weather"
	"github.com/DaniilKalts/rbk-school/6-week/internal/service/auth"
)

func WriteServiceError(w http.ResponseWriter, err error) {
	status, msg := http.StatusInternalServerError, "internal server error"
	switch {
	case errors.Is(err, user.ErrNotFound), errors.Is(err, city.ErrNotFound):
		status, msg = http.StatusNotFound, err.Error()
	case errors.Is(err, user.ErrEmailAlreadyExists), errors.Is(err, city.ErrAlreadyExists):
		status, msg = http.StatusConflict, err.Error()
	case errors.Is(err, user.ErrInvalidID),
		errors.Is(err, user.ErrInvalidFirstName),
		errors.Is(err, user.ErrInvalidLastName),
		errors.Is(err, user.ErrInvalidEmail),
		errors.Is(err, user.ErrInvalidPassword),
		errors.Is(err, user.ErrInvalidRole),
		errors.Is(err, city.ErrInvalidID),
		errors.Is(err, city.ErrInvalidUserID),
		errors.Is(err, city.ErrInvalidName),
		errors.Is(err, weather.ErrInvalidCity),
		errors.Is(err, weather.ErrInvalidLimit),
		errors.Is(err, weather.ErrInvalidOffset):
		status, msg = http.StatusBadRequest, err.Error()
	case errors.Is(err, auth.ErrInvalidCredentials), errors.Is(err, jwt.ErrInvalidToken):
		status, msg = http.StatusUnauthorized, err.Error()
	}
	WriteError(w, status, msg)
}
