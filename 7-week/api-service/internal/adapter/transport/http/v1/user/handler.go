package user

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/7-week/pkg/httpx"
	"github.com/DaniilKalts/rbk-school/7-week/api-service/internal/domain/user"
	"github.com/DaniilKalts/rbk-school/7-week/pkg/jwt"

	serviceuser "github.com/DaniilKalts/rbk-school/7-week/api-service/internal/service/user"
)

type Service interface {
	Create(ctx context.Context, input serviceuser.CreateInput) (*user.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*user.User, error)
	List(ctx context.Context) ([]user.User, error)
	Update(ctx context.Context, id uuid.UUID, input serviceuser.UpdateInput) (*user.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type TokenRevoker interface {
	Revoke(ctx context.Context, token string) error
}

type Handler struct {
	service      Service
	tokenRevoker TokenRevoker
}

func NewHandler(service Service, tokenRevoker TokenRevoker) *Handler {
	return &Handler{service: service, tokenRevoker: tokenRevoker}
}

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := httpx.CurrentUserID(w, r)
	if !ok {
		return
	}

	found, err := h.service.GetByID(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			httpx.WriteError(w, http.StatusNotFound, err.Error())
		default:
			httpx.WriteInternalError(w, r, err)
		}
		return
	}

	httpx.JSON(w, http.StatusOK, ToUserResponse(*found))
}

func (h *Handler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := httpx.CurrentUserID(w, r)
	if !ok {
		return
	}

	var body UpdateUserRequest
	if !httpx.DecodeJSON(w, r, &body) {
		return
	}

	updated, err := h.service.Update(r.Context(), userID, ToUpdateInput(body))
	if err != nil {
		switch {
		case errors.Is(err, user.ErrInvalidFirstName),
			errors.Is(err, user.ErrInvalidLastName),
			errors.Is(err, user.ErrInvalidEmail),
			errors.Is(err, user.ErrInvalidPassword),
			errors.Is(err, user.ErrInvalidRole):
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, user.ErrNotFound):
			httpx.WriteError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, user.ErrEmailAlreadyExists):
			httpx.WriteError(w, http.StatusConflict, err.Error())
		default:
			httpx.WriteInternalError(w, r, err)
		}
		return
	}

	httpx.JSON(w, http.StatusOK, ToUserResponse(*updated))
}

func (h *Handler) DeleteMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := httpx.CurrentUserID(w, r)
	if !ok {
		return
	}

	token, ok := httpx.BearerTokenFromRequest(r)
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "отсутствует или некорректный заголовок Authorization")
		return
	}

	if err := h.tokenRevoker.Revoke(r.Context(), token); err != nil {
		switch {
		case errors.Is(err, jwt.ErrInvalidToken):
			httpx.WriteError(w, http.StatusUnauthorized, err.Error())
		default:
			httpx.WriteInternalError(w, r, err)
		}
		return
	}

	if err := h.service.Delete(r.Context(), userID); err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			httpx.WriteError(w, http.StatusNotFound, err.Error())
		default:
			httpx.WriteInternalError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var body CreateUserRequest
	if !httpx.DecodeJSON(w, r, &body) {
		return
	}

	created, err := h.service.Create(r.Context(), ToCreateInput(body))
	if err != nil {
		switch {
		case errors.Is(err, user.ErrInvalidFirstName),
			errors.Is(err, user.ErrInvalidLastName),
			errors.Is(err, user.ErrInvalidEmail),
			errors.Is(err, user.ErrInvalidPassword),
			errors.Is(err, user.ErrInvalidRole):
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, user.ErrEmailAlreadyExists):
			httpx.WriteError(w, http.StatusConflict, err.Error())
		default:
			httpx.WriteInternalError(w, r, err)
		}
		return
	}

	httpx.JSON(w, http.StatusCreated, ToUserResponse(*created))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.List(r.Context())
	if err != nil {
		httpx.WriteInternalError(w, r, err)
		return
	}

	httpx.JSON(w, http.StatusOK, ToUserResponses(users))
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "некорректный id пользователя")
		return
	}

	found, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			httpx.WriteError(w, http.StatusNotFound, err.Error())
		default:
			httpx.WriteInternalError(w, r, err)
		}
		return
	}

	httpx.JSON(w, http.StatusOK, ToUserResponse(*found))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "некорректный id пользователя")
		return
	}

	var body UpdateUserRequest
	if !httpx.DecodeJSON(w, r, &body) {
		return
	}

	updated, err := h.service.Update(r.Context(), id, ToUpdateInput(body))
	if err != nil {
		switch {
		case errors.Is(err, user.ErrInvalidFirstName),
			errors.Is(err, user.ErrInvalidLastName),
			errors.Is(err, user.ErrInvalidEmail),
			errors.Is(err, user.ErrInvalidPassword),
			errors.Is(err, user.ErrInvalidRole):
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, user.ErrNotFound):
			httpx.WriteError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, user.ErrEmailAlreadyExists):
			httpx.WriteError(w, http.StatusConflict, err.Error())
		default:
			httpx.WriteInternalError(w, r, err)
		}
		return
	}

	httpx.JSON(w, http.StatusOK, ToUserResponse(*updated))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "некорректный id пользователя")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			httpx.WriteError(w, http.StatusNotFound, err.Error())
		default:
			httpx.WriteInternalError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
