package city

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/5-week/internal/adapters/transport/http/helpers"
	"github.com/DaniilKalts/rbk-school/5-week/internal/adapters/transport/http/v1/city/dto"
)

func (h *Handler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		userID, ok := CurrentUserID(w, req)
		if !ok || userID == uuid.Nil {
			helpers.JSON(w, http.StatusUnauthorized, helpers.NewErrorResponse(http.StatusUnauthorized, "отсутствуют claims аутентификации"))
			return
		}

		cities, err := h.service.List(req.Context(), userID)
		if err != nil {
			WriteServiceError(w, err)
			return
		}

		helpers.JSON(w, http.StatusOK, dto.ToCityResponses(cities))
	}
}
