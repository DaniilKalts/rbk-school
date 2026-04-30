package city

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/3-week/internal/adapters/transport/http/v1/city/dto"
	domaincity "github.com/DaniilKalts/rbk-school/3-week/internal/domain/city"
	domainuser "github.com/DaniilKalts/rbk-school/3-week/internal/domain/user"
	servicecity "github.com/DaniilKalts/rbk-school/3-week/internal/service/city"
	"github.com/DaniilKalts/rbk-school/3-week/internal/utils"
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
	userID, ok := parseUUID(w, r.PathValue("id"), "invalid user id")
	if !ok {
		return
	}

	var req dto.CreateCityRequest
	if !utils.DecodeJSON(w, r, &req) {
		return
	}

	c, err := h.service.Create(r.Context(), userID, dto.ToCreateInput(req))
	if err != nil {
		writeServiceError(w, err)
		return
	}

	utils.JSON(w, http.StatusCreated, dto.ToCityResponse(*c))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := parseUUID(w, r.PathValue("id"), "invalid user id")
	if !ok {
		return
	}

	cities, err := h.service.List(r.Context(), userID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	utils.JSON(w, http.StatusOK, dto.ToCityResponses(cities))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := parseUUID(w, r.PathValue("id"), "invalid user id")
	if !ok {
		return
	}

	cityID, ok := parseUUID(w, r.PathValue("city_id"), "invalid city id")
	if !ok {
		return
	}

	if err := h.service.Delete(r.Context(), userID, cityID); err != nil {
		writeServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseUUID(w http.ResponseWriter, value string, message string) (uuid.UUID, bool) {
	id, err := uuid.Parse(value)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, message)
		return uuid.Nil, false
	}

	return id, true
}
func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domainuser.ErrNotFound):
		utils.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domaincity.ErrNotFound):
		utils.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domaincity.ErrAlreadyExists):
		utils.Error(w, http.StatusConflict, err.Error())
	case errors.Is(err, domaincity.ErrInvalidID),
		errors.Is(err, domaincity.ErrInvalidUserID),
		errors.Is(err, domaincity.ErrInvalidName),
		errors.Is(err, domainuser.ErrInvalidID):
		utils.Error(w, http.StatusBadRequest, err.Error())
	default:
		utils.Error(w, http.StatusInternalServerError, "internal server error")
	}
}
