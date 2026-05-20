package weather

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/7-week/api-service/internal/adapter/database/postgres/sqlc"
	"github.com/DaniilKalts/rbk-school/7-week/api-service/internal/domain/history"
)

type Repository struct {
	queries *sqlc.Queries
}

func NewRepository(db sqlc.DBTX) *Repository {
	return &Repository{queries: sqlc.New(db)}
}

func (r *Repository) CreateHistory(ctx context.Context, h history.History) (*history.History, error) {
	row, err := r.queries.CreateWeatherHistory(ctx, sqlc.CreateWeatherHistoryParams{
		ID:          h.ID,
		UserID:      h.UserID,
		City:        h.City,
		Temperature: h.Temperature,
		Description: h.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("создание записи истории погоды: %w", err)
	}

	return new(toDomain(row)), nil
}

func (r *Repository) ListHistory(ctx context.Context, userID uuid.UUID, city string, limit int, offset int) ([]history.History, error) {
	rows, err := r.queries.ListWeatherHistory(ctx, sqlc.ListWeatherHistoryParams{
		UserID: userID,
		City:   city,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("получение истории погоды: %w", err)
	}

	history := make([]history.History, 0, len(rows))
	for _, row := range rows {
		history = append(history, toDomain(row))
	}

	return history, nil
}
