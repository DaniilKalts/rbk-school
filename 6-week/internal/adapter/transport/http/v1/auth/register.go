package auth

import (
	"net/http"

	"github.com/DaniilKalts/rbk-school/6-week/internal/adapter/transport/http/httpx"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var body RegisterRequest
	if !httpx.DecodeJSON(w, r, &body) {
		return
	}

	token, err := h.service.Register(r.Context(), ToRegisterInput(body))
	if err != nil {
		httpx.WriteServiceError(w, err)
		return
	}

	httpx.JSON(w, http.StatusCreated, ToTokenResponse(*token))
}
