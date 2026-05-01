package auth

import (
	"errors"
	"net/http"

	"github.com/DaniilKalts/rbk-school/4-week/internal/adapters/transport/http/helpers"
	domainuser "github.com/DaniilKalts/rbk-school/4-week/internal/domain/user"
	serviceauth "github.com/DaniilKalts/rbk-school/4-week/internal/service/auth"
)

func writeServiceError(w http.ResponseWriter, err error) {
	status, message := serviceErrorResponse(err)
	response := helpers.NewErrorResponse(status, message)
	helpers.JSON(w, status, response)
}

func serviceErrorResponse(err error) (int, string) {
	switch {
	case errors.Is(err, serviceauth.ErrInvalidCredentials):
		return http.StatusUnauthorized, err.Error()
	case errors.Is(err, domainuser.ErrEmailAlreadyExists):
		return http.StatusConflict, err.Error()
	case errors.Is(err, domainuser.ErrInvalidID),
		errors.Is(err, domainuser.ErrInvalidFirstName),
		errors.Is(err, domainuser.ErrInvalidLastName),
		errors.Is(err, domainuser.ErrInvalidEmail),
		errors.Is(err, domainuser.ErrInvalidPassword),
		errors.Is(err, domainuser.ErrInvalidRole):
		return http.StatusBadRequest, err.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
