package weather

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/DaniilKalts/rbk-school/5-week/internal/adapter/transport/http/httpx"

	domaincity "github.com/DaniilKalts/rbk-school/5-week/internal/domain/city"
)

func (h *Handler) History(w http.ResponseWriter, r *http.Request) {
	userID, ok := httpx.CurrentUserID(w, r)
	if !ok {
		return
	}

	city := strings.TrimSpace(r.URL.Query().Get("city"))
	limit, ok := parseLimit(w, r.URL.Query().Get("limit"))
	if !ok {
		return
	}
	offset, ok := parseOffset(w, r.URL.Query().Get("offset"))
	if !ok {
		return
	}

	history, err := h.service.GetHistory(r.Context(), userID, city, limit, offset)
	if err != nil {
		httpx.WriteServiceError(w, err)
		return
	}

	httpx.JSON(w, http.StatusOK, ToUserWeatherHistoryResponse(userID, domaincity.NormalizeCityName(city), history))
}

func parseLimit(w http.ResponseWriter, value string) (int, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, true
	}
	limit, err := strconv.Atoi(value)
	if err != nil || limit <= 0 || limit > math.MaxInt32 {
		httpx.WriteError(w, http.StatusBadRequest, "limit должен быть положительным числом")
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
		httpx.WriteError(w, http.StatusBadRequest, "offset должен быть неотрицательным числом")
		return 0, false
	}
	return offset, true
}
