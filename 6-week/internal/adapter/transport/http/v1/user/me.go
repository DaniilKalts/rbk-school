package user

import (
	"errors"
	"net/http"

	"github.com/DaniilKalts/rbk-school/6-week/internal/adapter/transport/http/httpx"
	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/user"
	"github.com/DaniilKalts/rbk-school/6-week/pkg/jwt"
)

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := httpx.CurrentUserID(w, r)
	if !ok {
		return
	}

	found, err := h.service.GetByID(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			httpx.WriteError(w, http.StatusNotFound, err.Error())
		default:
			httpx.WriteInternalError(w, r, err)
		}
		return
	}

	httpx.JSON(w, http.StatusOK, ToUserResponse(*found))
}

func (h *Handler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := httpx.CurrentUserID(w, r)
	if !ok {
		return
	}

	var body UpdateUserRequest
	if !httpx.DecodeJSON(w, r, &body) {
		return
	}

	updated, err := h.service.Update(r.Context(), userID, ToUpdateInput(body))
	if err != nil {
		switch {
		case errors.Is(err, user.ErrInvalidFirstName),
			errors.Is(err, user.ErrInvalidLastName),
			errors.Is(err, user.ErrInvalidEmail),
			errors.Is(err, user.ErrInvalidPassword),
			errors.Is(err, user.ErrInvalidRole):
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, user.ErrNotFound):
			httpx.WriteError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, user.ErrEmailAlreadyExists):
			httpx.WriteError(w, http.StatusConflict, err.Error())
		default:
			httpx.WriteInternalError(w, r, err)
		}
		return
	}

	httpx.JSON(w, http.StatusOK, ToUserResponse(*updated))
}

func (h *Handler) DeleteMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := httpx.CurrentUserID(w, r)
	if !ok {
		return
	}

	token, ok := httpx.BearerTokenFromRequest(r)
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "отсутствует или некорректный заголовок Authorization")
		return
	}

	if err := h.tokenRevoker.Revoke(r.Context(), token); err != nil {
		switch {
		case errors.Is(err, jwt.ErrInvalidToken):
			httpx.WriteError(w, http.StatusUnauthorized, err.Error())
		default:
			httpx.WriteInternalError(w, r, err)
		}
		return
	}

	if err := h.service.Delete(r.Context(), userID); err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			httpx.WriteError(w, http.StatusNotFound, err.Error())
		default:
			httpx.WriteInternalError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
