package weather

import (
	"net/http"

	"github.com/DaniilKalts/rbk-school/5-week/internal/adapters/transport/http/helpers"
	"github.com/DaniilKalts/rbk-school/5-week/internal/adapters/transport/http/v1/weather/dto"
)

func (h *Handler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		userID, ok := CurrentUserID(w, req)
		if !ok {
			return
		}

		weathers, err := h.service.GetByUserID(req.Context(), userID)
		if err != nil {
			WriteServiceError(w, err)
			return
		}

		helpers.JSON(w, http.StatusOK, dto.ToUserWeatherResponse(userID, weathers))
	}
}
