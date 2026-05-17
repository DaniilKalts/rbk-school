package city

import (
	"errors"
	"net/http"

	"github.com/DaniilKalts/rbk-school/6-week/internal/adapter/transport/http/httpx"
	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/city"
	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/user"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := httpx.CurrentUserID(w, r)
	if !ok {
		return
	}

	var body CreateCityRequest
	if !httpx.DecodeJSON(w, r, &body) {
		return
	}

	created, err := h.service.Create(r.Context(), userID, ToCreateInput(body))
	if err != nil {
		switch {
		case errors.Is(err, city.ErrInvalidName), errors.Is(err, city.ErrInvalidUserID):
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, user.ErrNotFound):
			httpx.WriteError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, city.ErrAlreadyExists):
			httpx.WriteError(w, http.StatusConflict, err.Error())
		default:
			httpx.WriteInternalError(w, r, err)
		}
		return
	}

	httpx.JSON(w, http.StatusCreated, ToCityResponse(*created))
}
