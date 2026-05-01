package helpers

import (
	"encoding/json"
	"net/http"
)

func DecodeJSON(w http.ResponseWriter, r *http.Request, v any) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		response := NewErrorResponse(http.StatusBadRequest, "некорректное тело запроса")
		JSON(w, http.StatusBadRequest, response)
		return false
	}

	return true
}
