package city

import (
	"errors"
	"net/http"

	"github.com/DaniilKalts/rbk-school/6-week/internal/adapter/transport/http/httpx"
	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/user"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := httpx.CurrentUserID(w, r)
	if !ok {
		return
	}

	cities, err := h.service.List(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			httpx.WriteError(w, http.StatusNotFound, err.Error())
		default:
			httpx.WriteInternalError(w, r, err)
		}
		return
	}

	httpx.JSON(w, http.StatusOK, ToCityResponses(cities))
}
