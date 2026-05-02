package city

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/4-week/internal/adapters/transport/http/helpers"
)

func (h *Handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		userID, ok := CurrentUserID(w, req)
		if !ok || userID == uuid.Nil {
			helpers.JSON(w, http.StatusUnauthorized, helpers.NewErrorResponse(http.StatusUnauthorized, "отсутствуют claims аутентификации"))
			return
		}

		cityID, err := uuid.Parse(chi.URLParam(req, "city_id"))
		if err != nil {
			helpers.JSON(w, http.StatusBadRequest, helpers.NewErrorResponse(http.StatusBadRequest, "invalid city id"))
			return
		}

		if err := h.service.Delete(req.Context(), userID, cityID); err != nil {
			WriteServiceError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
