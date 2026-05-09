package helpers

import (
	"encoding/json"
	"errors"
	"net/http"
)

const maxRequestBodyBytes = 1 << 20

func DecodeJSON(w http.ResponseWriter, r *http.Request, v any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(v); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			JSON(w, http.StatusRequestEntityTooLarge, NewErrorResponse(http.StatusRequestEntityTooLarge, "тело запроса слишком большое"))
			return false
		}

		JSON(w, http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "некорректное тело запроса"))
		return false
	}

	return true
}
