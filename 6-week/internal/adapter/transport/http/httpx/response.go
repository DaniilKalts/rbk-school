package httpx

import (
	"encoding/json"
	"net/http"

	"github.com/DaniilKalts/rbk-school/6-week/pkg/logger"
	"go.uber.org/zap"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, message string) {
	JSON(w, status, ErrorResponse{Code: status, Message: message})
}

func WriteInternalError(w http.ResponseWriter, r *http.Request, err error) {
	logger.FromContext(r.Context()).Error("необработанная ошибка сервиса", zap.Error(err))
	WriteError(w, http.StatusInternalServerError, "internal server error")
}
