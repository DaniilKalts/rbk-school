package auth

import (
	"net/http"

	"github.com/DaniilKalts/rbk-school/4-week/internal/adapters/transport/http/helpers"
)

func (h *Handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		token, ok := helpers.BearerTokenFromRequest(req)
		if !ok {
			helpers.JSON(w, http.StatusUnauthorized, helpers.NewErrorResponse(http.StatusUnauthorized, "отсутствует или некорректный заголовок Authorization"))
			return
		}

		if err := h.service.Logout(req.Context(), token); err != nil {
			helpers.JSON(w, http.StatusUnauthorized, helpers.NewErrorResponse(http.StatusUnauthorized, "некорректный или просроченный токен"))
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
