package auth

import (
	"net/http"

	"github.com/DaniilKalts/rbk-school/5-week/internal/adapters/transport/http/helpers"
	"github.com/DaniilKalts/rbk-school/5-week/internal/adapters/transport/http/v1/auth/dto"
)

func (h *Handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var body dto.LoginRequest
		if !helpers.DecodeJSON(w, req, &body) {
			return
		}

		token, err := h.service.Login(req.Context(), dto.ToLoginInput(body))
		if err != nil {
			WriteServiceError(w, err)
			return
		}

		helpers.JSON(w, http.StatusOK, dto.ToTokenResponse(*token))
	}
}
