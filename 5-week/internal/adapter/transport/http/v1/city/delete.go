package city

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/5-week/internal/adapter/transport/http/httpx"
)

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := httpx.CurrentUserID(w, r)
	if !ok {
		return
	}

	cityID, err := uuid.Parse(chi.URLParam(r, "city_id"))
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid city id")
		return
	}

	if err := h.service.Delete(r.Context(), userID, cityID); err != nil {
		httpx.WriteServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
