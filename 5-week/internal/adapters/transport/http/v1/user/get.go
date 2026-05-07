package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/5-week/internal/adapters/transport/http/helpers"
	"github.com/DaniilKalts/rbk-school/5-week/internal/adapters/transport/http/v1/user/dto"
)

func (h *Handler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		users, err := h.service.List(req.Context())
		if err != nil {
			WriteServiceError(w, err)
			return
		}

		helpers.JSON(w, http.StatusOK, dto.ToUserResponses(users))
	}
}

func (h *Handler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id, err := uuid.Parse(chi.URLParam(req, "id"))
		if err != nil {
			helpers.JSON(w, http.StatusBadRequest, helpers.NewErrorResponse(http.StatusBadRequest, "некорректный id пользователя"))
			return
		}

		found, err := h.service.GetByID(req.Context(), id)
		if err != nil {
			WriteServiceError(w, err)
			return
		}

		helpers.JSON(w, http.StatusOK, dto.ToUserResponse(*found))
	}
}

func (h *Handler) Me() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		userID, ok := currentUserID(w, req)
		if !ok {
			return
		}

		found, err := h.service.GetByID(req.Context(), userID)
		if err != nil {
			WriteServiceError(w, err)
			return
		}

		helpers.JSON(w, http.StatusOK, dto.ToUserResponse(*found))
	}
}
