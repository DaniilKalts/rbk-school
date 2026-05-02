package city

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/4-week/internal/adapters/transport/http/helpers"
	"github.com/DaniilKalts/rbk-school/4-week/internal/adapters/transport/http/v1/city/dto"
)

func (h *Handler) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		userID, ok := CurrentUserID(w, req)
		if !ok || userID == uuid.Nil {
			helpers.JSON(w, http.StatusUnauthorized, helpers.NewErrorResponse(http.StatusUnauthorized, "отсутствуют claims аутентификации"))
			return
		}

		var body dto.CreateCityRequest
		if !helpers.DecodeJSON(w, req, &body) {
			return
		}

		created, err := h.service.Create(req.Context(), userID, dto.ToCreateInput(body))
		if err != nil {
			WriteServiceError(w, err)
			return
		}

		helpers.JSON(w, http.StatusCreated, dto.ToCityResponse(*created))
	}
}
