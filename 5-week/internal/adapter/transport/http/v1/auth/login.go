package auth

import (
	"net/http"

	"github.com/DaniilKalts/rbk-school/5-week/internal/adapter/transport/http/httpx"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var body LoginRequest
	if !httpx.DecodeJSON(w, r, &body) {
		return
	}

	token, err := h.service.Login(r.Context(), ToLoginInput(body))
	if err != nil {
		httpx.WriteServiceError(w, err)
		return
	}

	httpx.JSON(w, http.StatusOK, ToTokenResponse(*token))
}
