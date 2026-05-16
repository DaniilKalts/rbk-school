package city

import (
	"net/http"

	"github.com/DaniilKalts/rbk-school/5-week/internal/adapter/transport/http/httpx"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := httpx.CurrentUserID(w, r)
	if !ok {
		return
	}

	var body CreateCityRequest
	if !httpx.DecodeJSON(w, r, &body) {
		return
	}

	created, err := h.service.Create(r.Context(), userID, ToCreateInput(body))
	if err != nil {
		httpx.WriteServiceError(w, err)
		return
	}

	httpx.JSON(w, http.StatusCreated, ToCityResponse(*created))
}
