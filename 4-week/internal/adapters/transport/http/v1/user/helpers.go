package user

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/4-week/internal/adapters/transport/http/helpers"
	domainuser "github.com/DaniilKalts/rbk-school/4-week/internal/domain/user"
)

func currentUserID(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	claims, ok := helpers.ClaimsFromContext(r.Context())
	if !ok {
		response := helpers.NewErrorResponse(http.StatusUnauthorized, "отсутствуют claims аутентификации")
		helpers.JSON(w, http.StatusUnauthorized, response)
		return uuid.Nil, false
	}

	return claims.UserID, true
}

func parseID(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response := helpers.NewErrorResponse(http.StatusBadRequest, "некорректный id пользователя")
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
