package auth

import (
	"net/http"

	"github.com/DaniilKalts/rbk-school/5-week/internal/adapters/transport/http/helpers"
)

func (h *Handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		token, ok := helpers.BearerTokenFromRequest(req)
		if !ok {
			helpers.JSON(w, http.StatusUnauthorized, helpers.NewErrorResponse(http.StatusUnauthorized, "отсутствует или некорректный заголовок Authorization"))
			return
		}

		if err := h.service.Logout(req.Context(), token); err != nil {
			WriteServiceError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
