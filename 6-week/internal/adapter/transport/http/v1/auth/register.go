package auth

import (
	"errors"
	"net/http"

	"github.com/DaniilKalts/rbk-school/6-week/internal/adapter/transport/http/httpx"
	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/user"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var body RegisterRequest
	if !httpx.DecodeJSON(w, r, &body) {
		return
	}

	token, err := h.service.Register(r.Context(), ToRegisterInput(body))
	if err != nil {
		switch {
		case errors.Is(err, user.ErrInvalidFirstName),
			errors.Is(err, user.ErrInvalidLastName),
			errors.Is(err, user.ErrInvalidEmail),
			errors.Is(err, user.ErrInvalidPassword):
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, user.ErrEmailAlreadyExists):
			httpx.WriteError(w, http.StatusConflict, err.Error())
		default:
			httpx.WriteInternalError(w, r, err)
		}
		return
	}

	httpx.JSON(w, http.StatusCreated, ToTokenResponse(*token))
}
