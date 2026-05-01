package weather

import (
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/4-week/internal/adapters/transport/http/helpers"
	domainuser "github.com/DaniilKalts/rbk-school/4-week/internal/domain/user"
	domainweather "github.com/DaniilKalts/rbk-school/4-week/internal/domain/weather"
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

func parseLimit(w http.ResponseWriter, value string) (int, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, true
	}

	limit, err := strconv.Atoi(value)
	if err != nil || limit <= 0 || limit > math.MaxInt32 {
		response := helpers.NewErrorResponse(http.StatusBadRequest, "limit должен быть положительным числом")
		helpers.JSON(w, http.StatusBadRequest, response)
		return 0, false
	}

	return limit, true
}

func parseOffset(w http.ResponseWriter, value string) (int, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, true
	}

	offset, err := strconv.Atoi(value)
	if err != nil || offset < 0 || offset > math.MaxInt32 {
		response := helpers.NewErrorResponse(http.StatusBadRequest, "offset должен быть неотрицательным числом")
		helpers.JSON(w, http.StatusBadRequest, response)
		return 0, false
	}

	return offset, true
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
	case errors.Is(err, domainuser.ErrInvalidID),
		errors.Is(err, domainweather.ErrInvalidCity),
		errors.Is(err, domainweather.ErrInvalidLimit),
		errors.Is(err, domainweather.ErrInvalidOffset):
		return http.StatusBadRequest, err.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
