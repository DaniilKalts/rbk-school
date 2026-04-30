package user

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/3-week/internal/adapters/transport/http/v1/user/dto"
	domainuser "github.com/DaniilKalts/rbk-school/3-week/internal/domain/user"
	serviceuser "github.com/DaniilKalts/rbk-school/3-week/internal/service/user"
	"github.com/DaniilKalts/rbk-school/3-week/internal/utils"
)

type Service interface {
	Create(ctx context.Context, input serviceuser.CreateInput) (*domainuser.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domainuser.User, error)
	List(ctx context.Context) ([]domainuser.User, error)
	Update(ctx context.Context, id uuid.UUID, input serviceuser.UpdateInput) (*domainuser.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if !utils.DecodeJSON(w, r, &req) {
		return
	}

	u, err := h.service.Create(r.Context(), dto.ToCreateInput(req))
	if err != nil {
		writeServiceError(w, err)
		return
	}

	utils.JSON(w, http.StatusCreated, dto.ToUserResponse(*u))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.List(r.Context())
	if err != nil {
		writeServiceError(w, err)
		return
	}

	utils.JSON(w, http.StatusOK, dto.ToUserResponses(users))
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	u, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	utils.JSON(w, http.StatusOK, dto.ToUserResponse(*u))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	var req dto.UpdateUserRequest
	if !utils.DecodeJSON(w, r, &req) {
		return
	}

	u, err := h.service.Update(r.Context(), id, dto.ToUpdateInput(req))
	if err != nil {
		writeServiceError(w, err)
		return
	}

	utils.JSON(w, http.StatusOK, dto.ToUserResponse(*u))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		writeServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseID(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid user id")
		return uuid.Nil, false
	}

	return id, true
}
func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domainuser.ErrNotFound):
		utils.Error(w, http.StatusNotFound, err.Error())
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
