package city

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/4-week/internal/adapters/transport/http/helpers"
	domaincity "github.com/DaniilKalts/rbk-school/4-week/internal/domain/city"
	domainuser "github.com/DaniilKalts/rbk-school/4-week/internal/domain/user"
)

func parseUUID(w http.ResponseWriter, value string, message string) (uuid.UUID, bool) {
	id, err := uuid.Parse(value)
	if err != nil {
		response := helpers.NewErrorResponse(http.StatusBadRequest, message)
		helpers.JSON(w, http.StatusBadRequest, response)
		return uuid.Nil, false
	}

	return id, true
}

func writeServiceError(w http.ResponseWriter, err error) {
	status, message := serviceErrorResponse(err)
	response := helpers.NewErrorResponse(status, message)
	helpers.JSON(w, status, response)
}

func serviceErrorResponse(err error) (int, string) {
	switch {
	case errors.Is(err, domainuser.ErrNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, domaincity.ErrNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, domaincity.ErrAlreadyExists):
		return http.StatusConflict, err.Error()
	case errors.Is(err, domaincity.ErrInvalidID),
		errors.Is(err, domaincity.ErrInvalidUserID),
		errors.Is(err, domaincity.ErrInvalidName),
		errors.Is(err, domainuser.ErrInvalidID):
		return http.StatusBadRequest, err.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
