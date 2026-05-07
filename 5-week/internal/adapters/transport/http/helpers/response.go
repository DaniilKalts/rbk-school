package helpers

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewErrorResponse(status int, message string) ErrorResponse {
	return ErrorResponse{Code: status, Message: message}
}

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
