package user

import (
	"net/http"

	"github.com/DaniilKalts/rbk-school/4-week/internal/adapters/transport/http/helpers"
	"github.com/DaniilKalts/rbk-school/4-week/internal/adapters/transport/http/v1/user/dto"
)

func (h *Handler) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var body dto.CreateUserRequest
		if !helpers.DecodeJSON(w, req, &body) {
			return
		}

		created, err := h.service.Create(req.Context(), dto.ToCreateInput(body))
		if err != nil {
			WriteServiceError(w, err)
			return
		}

		helpers.JSON(w, http.StatusCreated, dto.ToUserResponse(*created))
	}
}
