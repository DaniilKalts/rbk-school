package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/5-week/internal/adapters/transport/http/helpers"
	"github.com/DaniilKalts/rbk-school/5-week/internal/adapters/transport/http/v1/user/dto"
)

func (h *Handler) Patch() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id, err := uuid.Parse(chi.URLParam(req, "id"))
		if err != nil {
			helpers.JSON(w, http.StatusBadRequest, helpers.NewErrorResponse(http.StatusBadRequest, "некорректный id пользователя"))
			return
		}

		var body dto.UpdateUserRequest
		if !helpers.DecodeJSON(w, req, &body) {
			return
		}

		updated, err := h.service.Update(req.Context(), id, dto.ToUpdateInput(body))
		if err != nil {
			WriteServiceError(w, err)
			return
		}

		helpers.JSON(w, http.StatusOK, dto.ToUserResponse(*updated))
	}
}

func (h *Handler) MePatch() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		userID, ok := currentUserID(w, req)
		if !ok {
			return
		}

		var body dto.UpdateUserRequest
		if !helpers.DecodeJSON(w, req, &body) {
			return
		}

		updated, err := h.service.Update(req.Context(), userID, dto.ToUpdateInput(body))
		if err != nil {
			WriteServiceError(w, err)
			return
		}

		helpers.JSON(w, http.StatusOK, dto.ToUserResponse(*updated))
	}
}
