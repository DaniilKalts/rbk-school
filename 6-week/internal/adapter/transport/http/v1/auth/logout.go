package auth

import (
	"net/http"

	"github.com/DaniilKalts/rbk-school/6-week/internal/adapter/transport/http/httpx"
)

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	token, ok := httpx.BearerTokenFromRequest(r)
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "отсутствует или некорректный заголовок Authorization")
		return
	}

	if err := h.service.Logout(r.Context(), token); err != nil {
		httpx.WriteServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
