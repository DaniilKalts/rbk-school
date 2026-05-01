package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/DaniilKalts/rbk-school/4-week/internal/adapters/transport/http/v1/auth/dto"
	domainuser "github.com/DaniilKalts/rbk-school/4-week/internal/domain/user"
	serviceauth "github.com/DaniilKalts/rbk-school/4-week/internal/service/auth"
	"github.com/DaniilKalts/rbk-school/4-week/internal/utils"
)

type Service interface {
	Register(ctx context.Context, input serviceauth.RegisterInput) (*serviceauth.Token, error)
	Login(ctx context.Context, input serviceauth.LoginInput) (*serviceauth.Token, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if !utils.DecodeJSON(w, r, &req) {
		return
	}

	token, err := h.service.Register(r.Context(), dto.ToRegisterInput(req))
	if err != nil {
		writeServiceError(w, err)
		return
	}

	utils.JSON(w, http.StatusCreated, dto.ToTokenResponse(*token))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if !utils.DecodeJSON(w, r, &req) {
		return
	}

	token, err := h.service.Login(r.Context(), dto.ToLoginInput(req))
	if err != nil {
		writeServiceError(w, err)
		return
	}

	utils.JSON(w, http.StatusOK, dto.ToTokenResponse(*token))
}
func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, serviceauth.ErrInvalidCredentials):
		utils.Error(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, domainuser.ErrEmailAlreadyExists):
		utils.Error(w, http.StatusConflict, err.Error())
	case errors.Is(err, domainuser.ErrInvalidID),
		errors.Is(err, domainuser.ErrInvalidFirstName),
		errors.Is(err, domainuser.ErrInvalidLastName),
		errors.Is(err, domainuser.ErrInvalidEmail),
		errors.Is(err, domainuser.ErrInvalidPassword),
		errors.Is(err, domainuser.ErrInvalidRole):
		utils.Error(w, http.StatusBadRequest, err.Error())
	default:
		utils.Error(w, http.StatusInternalServerError, "internal server error")
	}
}
