package auth

import (
	"errors"
	"net/http"

	"github.com/DaniilKalts/rbk-school/6-week/internal/adapter/transport/http/httpx"
	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/user"
	"github.com/DaniilKalts/rbk-school/6-week/internal/service/auth"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var body LoginRequest
	if !httpx.DecodeJSON(w, r, &body) {
		return
	}

	token, err := h.service.Login(r.Context(), ToLoginInput(body))
	if err != nil {
		switch {
		case errors.Is(err, user.ErrInvalidEmail), errors.Is(err, user.ErrInvalidPassword):
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, auth.ErrInvalidCredentials):
			httpx.WriteError(w, http.StatusUnauthorized, err.Error())
		default:
			httpx.WriteInternalError(w, r, err)
		}
		return
	}

	httpx.JSON(w, http.StatusOK, ToTokenResponse(*token))
}
