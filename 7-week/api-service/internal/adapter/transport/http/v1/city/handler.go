package city

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/7-week/api-service/internal/adapter/transport/http/httpx"
	"github.com/DaniilKalts/rbk-school/7-week/api-service/internal/domain/city"
	"github.com/DaniilKalts/rbk-school/7-week/api-service/internal/domain/user"

	servicecity "github.com/DaniilKalts/rbk-school/7-week/api-service/internal/service/city"
)

type Service interface {
	Create(ctx context.Context, userID uuid.UUID, input servicecity.CreateInput) (*city.City, error)
	List(ctx context.Context, userID uuid.UUID) ([]city.City, error)
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := httpx.CurrentUserID(w, r)
	if !ok {
		return
	}

	var body CreateCityRequest
	if !httpx.DecodeJSON(w, r, &body) {
		return
	}

	created, err := h.service.Create(r.Context(), userID, ToCreateInput(body))
	if err != nil {
		switch {
		case errors.Is(err, city.ErrInvalidName), errors.Is(err, city.ErrInvalidUserID):
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, user.ErrNotFound):
			httpx.WriteError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, city.ErrAlreadyExists):
			httpx.WriteError(w, http.StatusConflict, err.Error())
		default:
			httpx.WriteInternalError(w, r, err)
		}
		return
	}

	httpx.JSON(w, http.StatusCreated, ToCityResponse(*created))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := httpx.CurrentUserID(w, r)
	if !ok {
		return
	}

	cities, err := h.service.List(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			httpx.WriteError(w, http.StatusNotFound, err.Error())
		default:
			httpx.WriteInternalError(w, r, err)
		}
		return
	}

	httpx.JSON(w, http.StatusOK, ToCityResponses(cities))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := httpx.CurrentUserID(w, r)
	if !ok {
		return
	}

	cityID, err := uuid.Parse(chi.URLParam(r, "city_id"))
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid city id")
		return
	}

	if err := h.service.Delete(r.Context(), userID, cityID); err != nil {
		switch {
		case errors.Is(err, city.ErrNotFound), errors.Is(err, user.ErrNotFound):
			httpx.WriteError(w, http.StatusNotFound, err.Error())
		default:
			httpx.WriteInternalError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
