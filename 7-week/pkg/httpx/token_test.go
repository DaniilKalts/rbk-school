package httpx_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DaniilKalts/rbk-school/7-week/pkg/httpx"
)

func TestBearerTokenFromRequest(t *testing.T) {
	tests := []struct {
		name      string
		header    string
		wantToken string
		wantOK    bool
	}{
		{name: "valid bearer", header: "Bearer abc.def.ghi", wantToken: "abc.def.ghi", wantOK: true},
		{name: "case-insensitive scheme", header: "bearer abc.def.ghi", wantToken: "abc.def.ghi", wantOK: true},
		{name: "empty header", header: "", wantOK: false},
		{name: "missing token", header: "Bearer", wantOK: false},
		{name: "wrong scheme", header: "Basic abc", wantOK: false},
		{name: "extra parts", header: "Bearer abc def", wantOK: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.header != "" {
				req.Header.Set("Authorization", tc.header)
			}

			token, ok := httpx.BearerTokenFromRequest(req)
			assert.Equal(t, tc.wantOK, ok)
			assert.Equal(t, tc.wantToken, token)
		})
	}
}
