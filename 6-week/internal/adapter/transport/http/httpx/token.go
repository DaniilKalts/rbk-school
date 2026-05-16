package httpx

import (
	"net/http"
	"strings"
)

func BearerTokenFromRequest(r *http.Request) (string, bool) {
	fields := strings.Fields(r.Header.Get("Authorization"))
	if len(fields) != 2 || !strings.EqualFold(fields[0], "Bearer") {
		return "", false
	}

	return fields[1], true
}
