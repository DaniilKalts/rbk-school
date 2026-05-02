package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/4-week/internal/adapters/transport/http/helpers"
)

func (h *Handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id, err := uuid.Parse(chi.URLParam(req, "id"))
		if err != nil {
			helpers.JSON(w, http.StatusBadRequest, helpers.NewErrorResponse(http.StatusBadRequest, "некорректный id пользователя"))
			return
		}

		if err := h.service.Delete(req.Context(), id); err != nil {
			WriteServiceError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *Handler) MeDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		userID, ok := currentUserID(w, req)
		if !ok {
			return
		}

		if err := h.service.Delete(req.Context(), userID); err != nil {
			WriteServiceError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
