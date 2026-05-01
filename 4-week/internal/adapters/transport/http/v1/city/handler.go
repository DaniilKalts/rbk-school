package city

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/4-week/internal/adapters/transport/http/helpers"
	"github.com/DaniilKalts/rbk-school/4-week/internal/adapters/transport/http/v1/city/dto"
	domaincity "github.com/DaniilKalts/rbk-school/4-week/internal/domain/city"
	servicecity "github.com/DaniilKalts/rbk-school/4-week/internal/service/city"
)

type Service interface {
	Create(ctx context.Context, userID uuid.UUID, input servicecity.CreateInput) (*domaincity.City, error)
	List(ctx context.Context, userID uuid.UUID) ([]domaincity.City, error)
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := currentUserID(w, r)
	if !ok {
		return
	}

	var req dto.CreateCityRequest
	if !helpers.DecodeJSON(w, r, &req) {
		return
	}

	c, err := h.service.Create(r.Context(), userID, dto.ToCreateInput(req))
	if err != nil {
		writeServiceError(w, err)
		return
	}

	helpers.JSON(w, http.StatusCreated, dto.ToCityResponse(*c))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := currentUserID(w, r)
	if !ok {
		return
	}

	cities, err := h.service.List(r.Context(), userID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	helpers.JSON(w, http.StatusOK, dto.ToCityResponses(cities))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := currentUserID(w, r)
	if !ok {
		return
	}

	cityID, ok := parseUUIDParam(w, r, "city_id", "invalid city id")
	if !ok {
		return
	}

	if err := h.service.Delete(r.Context(), userID, cityID); err != nil {
		writeServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
