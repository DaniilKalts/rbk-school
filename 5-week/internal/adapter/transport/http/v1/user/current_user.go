package user

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/5-week/internal/adapter/transport/http/helpers"
)

func currentUserID(w http.ResponseWriter, req *http.Request) (uuid.UUID, bool) {
	claims, ok := helpers.ClaimsFromContext(req.Context())
	if !ok || claims.UserID == uuid.Nil {
		helpers.JSON(w, http.StatusUnauthorized, helpers.NewErrorResponse(http.StatusUnauthorized, "отсутствуют claims аутентификации"))
		return uuid.Nil, false
	}

	return claims.UserID, true
}
