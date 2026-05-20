package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/DaniilKalts/rbk-school/7-week/api-service/internal/adapter/transport/http/httpx"
	"github.com/DaniilKalts/rbk-school/7-week/api-service/internal/domain/user"
	"github.com/DaniilKalts/rbk-school/7-week/api-service/internal/service/auth"
	"github.com/DaniilKalts/rbk-school/7-week/pkg/jwt"
)

type Service interface {
	Register(ctx context.Context, input auth.RegisterInput) (*auth.Token, error)
	Login(ctx context.Context, input auth.LoginInput) (*auth.Token, error)
	Logout(ctx context.Context, accessToken string) error
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var body RegisterRequest
	if !httpx.DecodeJSON(w, r, &body) {
		return
	}

	token, err := h.service.Register(r.Context(), ToRegisterInput(body))
	if err != nil {
		switch {
		case errors.Is(err, user.ErrInvalidFirstName),
			errors.Is(err, user.ErrInvalidLastName),
			errors.Is(err, user.ErrInvalidEmail),
			errors.Is(err, user.ErrInvalidPassword):
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, user.ErrEmailAlreadyExists):
			httpx.WriteError(w, http.StatusConflict, err.Error())
		default:
			httpx.WriteInternalError(w, r, err)
		}
		return
	}

	httpx.JSON(w, http.StatusCreated, ToTokenResponse(*token))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var body LoginRequest
	if !httpx.DecodeJSON(w, r, &body) {
		return
	}

	token, err := h.service.Login(r.Context(), ToLoginInput(body))
	if err != nil {
		switch {
		case errors.Is(err, user.ErrInvalidEmail), errors.Is(err, user.ErrInvalidPassword):
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, auth.ErrInvalidCredentials):
			httpx.WriteError(w, http.StatusUnauthorized, err.Error())
		default:
			httpx.WriteInternalError(w, r, err)
		}
		return
	}

	httpx.JSON(w, http.StatusOK, ToTokenResponse(*token))
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	token, ok := httpx.BearerTokenFromRequest(r)
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "отсутствует или некорректный заголовок Authorization")
		return
	}

	if err := h.service.Logout(r.Context(), token); err != nil {
		switch {
		case errors.Is(err, jwt.ErrInvalidToken):
			httpx.WriteError(w, http.StatusUnauthorized, err.Error())
		default:
			httpx.WriteInternalError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
