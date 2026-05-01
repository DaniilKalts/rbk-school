package auth

import (
	"context"
	"net/http"

	"github.com/DaniilKalts/rbk-school/4-week/internal/adapters/transport/http/helpers"
	"github.com/DaniilKalts/rbk-school/4-week/internal/adapters/transport/http/v1/auth/dto"
	serviceauth "github.com/DaniilKalts/rbk-school/4-week/internal/service/auth"
)

type Service interface {
	Register(ctx context.Context, input serviceauth.RegisterInput) (*serviceauth.Token, error)
	Login(ctx context.Context, input serviceauth.LoginInput) (*serviceauth.Token, error)
	Logout(ctx context.Context, accessToken string) error
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if !helpers.DecodeJSON(w, r, &req) {
		return
	}

	token, err := h.service.Register(r.Context(), dto.ToRegisterInput(req))
	if err != nil {
		writeServiceError(w, err)
		return
	}

	helpers.JSON(w, http.StatusCreated, dto.ToTokenResponse(*token))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if !helpers.DecodeJSON(w, r, &req) {
		return
	}

	token, err := h.service.Login(r.Context(), dto.ToLoginInput(req))
	if err != nil {
		writeServiceError(w, err)
		return
	}

	helpers.JSON(w, http.StatusOK, dto.ToTokenResponse(*token))
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	token, ok := helpers.BearerTokenFromRequest(r)
	if !ok {
		response := helpers.NewErrorResponse(http.StatusUnauthorized, "missing or malformed authorization header")
		helpers.JSON(w, http.StatusUnauthorized, response)
		return
	}

	err := h.service.Logout(r.Context(), token)
	if err != nil {
		response := helpers.NewErrorResponse(http.StatusUnauthorized, "invalid or expired token")
		helpers.JSON(w, http.StatusUnauthorized, response)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
