package auth

import (
	"net/http"

	"github.com/DaniilKalts/rbk-school/4-week/internal/adapters/transport/http/helpers"
	"github.com/DaniilKalts/rbk-school/4-week/internal/adapters/transport/http/v1/auth/dto"
)

func (h *Handler) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var body dto.RegisterRequest
		if !helpers.DecodeJSON(w, req, &body) {
			return
		}

		token, err := h.service.Register(req.Context(), dto.ToRegisterInput(body))
		if err != nil {
			WriteServiceError(w, err)
			return
		}

		helpers.JSON(w, http.StatusCreated, dto.ToTokenResponse(*token))
	}
}
