package user

import (
	"net/http"

	"github.com/DaniilKalts/rbk-school/6-week/internal/adapter/transport/http/httpx"
)

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := httpx.CurrentUserID(w, r)
	if !ok {
		return
	}

	found, err := h.service.GetByID(r.Context(), userID)
	if err != nil {
		httpx.WriteServiceError(w, err)
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
		httpx.WriteServiceError(w, err)
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
		httpx.WriteServiceError(w, err)
		return
	}

	if err := h.service.Delete(r.Context(), userID); err != nil {
		httpx.WriteServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
