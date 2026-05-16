package weather

import (
	"net/http"

	"github.com/DaniilKalts/rbk-school/5-week/internal/adapter/transport/http/httpx"
)

func (h *Handler) GetByUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := httpx.CurrentUserID(w, r)
	if !ok {
		return
	}

	weathers, err := h.service.GetByUserID(r.Context(), userID)
	if err != nil {
		httpx.WriteServiceError(w, err)
		return
	}

	httpx.JSON(w, http.StatusOK, ToUserWeatherResponse(userID, weathers))
}
