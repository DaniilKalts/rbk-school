package weather

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/DaniilKalts/rbk-school/5-week/internal/adapters/transport/http/helpers"
	"github.com/DaniilKalts/rbk-school/5-week/internal/adapters/transport/http/v1/weather/dto"

	domaincity "github.com/DaniilKalts/rbk-school/5-week/internal/domain/city"
)

func (h *Handler) History() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		userID, ok := CurrentUserID(w, req)
		if !ok {
			return
		}

		city := strings.TrimSpace(req.URL.Query().Get("city"))
		limit, ok := parseLimit(w, req.URL.Query().Get("limit"))
		if !ok {
			return
		}
		offset, ok := parseOffset(w, req.URL.Query().Get("offset"))
		if !ok {
			return
		}

		history, err := h.service.GetHistory(req.Context(), userID, city, limit, offset)
		if err != nil {
			WriteServiceError(w, err)
			return
		}

		helpers.JSON(w, http.StatusOK, dto.ToUserWeatherHistoryResponse(userID, domaincity.NormalizeCityName(city), history))
	}
}

func parseLimit(w http.ResponseWriter, value string) (int, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, true
	}
	limit, err := strconv.Atoi(value)
	if err != nil || limit <= 0 || limit > math.MaxInt32 {
		helpers.JSON(w, http.StatusBadRequest, helpers.NewErrorResponse(http.StatusBadRequest, "limit должен быть положительным числом"))
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
		helpers.JSON(w, http.StatusBadRequest, helpers.NewErrorResponse(http.StatusBadRequest, "offset должен быть неотрицательным числом"))
		return 0, false
	}
	return offset, true
}
