package httpx_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/DaniilKalts/rbk-school/7-week/internal/adapter/transport/http/httpx"
)

type sampleBody struct {
	Name string `json:"name"`
}

func TestDecodeJSON_Success(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"Almaty"}`))
	w := httptest.NewRecorder()

	var got sampleBody
	ok := httpx.DecodeJSON(w, req, &got)

	require.True(t, ok)
	assert.Equal(t, "Almaty", got.Name)
	assert.Equal(t, http.StatusOK, w.Code, "не должно записывать ответ при успехе")
}

func TestDecodeJSON_MalformedJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{not json`))
	w := httptest.NewRecorder()

	var got sampleBody
	ok := httpx.DecodeJSON(w, req, &got)

	require.False(t, ok)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "некорректное тело запроса")
}

func TestDecodeJSON_UnknownFieldsRejected(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"x","extra":1}`))
	w := httptest.NewRecorder()

	var got sampleBody
	ok := httpx.DecodeJSON(w, req, &got)

	require.False(t, ok)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDecodeJSON_BodyTooLarge(t *testing.T) {
	huge := bytes.Repeat([]byte("a"), (1<<20)+10)
	body := bytes.NewBuffer(nil)
	body.WriteByte('"')
	body.Write(huge)
	body.WriteByte('"')

	req := httptest.NewRequest(http.MethodPost, "/", body)
	w := httptest.NewRecorder()

	var got string
	ok := httpx.DecodeJSON(w, req, &got)

	require.False(t, ok)
	assert.Equal(t, http.StatusRequestEntityTooLarge, w.Code)
	assert.Contains(t, w.Body.String(), "тело запроса слишком большое")
}
