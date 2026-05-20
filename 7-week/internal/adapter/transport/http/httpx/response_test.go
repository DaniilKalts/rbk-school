package httpx_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/DaniilKalts/rbk-school/7-week/internal/adapter/transport/http/httpx"
)

func TestJSON_WritesStatusAndBody(t *testing.T) {
	w := httptest.NewRecorder()
	httpx.JSON(w, http.StatusCreated, map[string]string{"hello": "world"})

	require.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"hello":"world"}`, w.Body.String())
}

func TestWriteError_ReturnsErrorEnvelope(t *testing.T) {
	w := httptest.NewRecorder()
	httpx.WriteError(w, http.StatusBadRequest, "плохие данные")

	require.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"code":400,"message":"плохие данные"}`, w.Body.String())
}

func TestWriteInternalError_Returns500WithGenericMessage(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	httpx.WriteInternalError(w, req, errors.New("boom"))

	require.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"code":500,"message":"internal server error"}`, w.Body.String())
}
