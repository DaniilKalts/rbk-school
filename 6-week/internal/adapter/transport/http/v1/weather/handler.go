package weather

import (
	"context"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/history"
	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/weather"
)

type Service interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]weather.Weather, error)
	GetHistory(ctx context.Context, userID uuid.UUID, city string, limit int, offset int) ([]history.History, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}
