package weather

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/4-week/internal/adapters/database/postgres/sqlc"
	domainhistory "github.com/DaniilKalts/rbk-school/4-week/internal/domain/history"
)

type Repository struct {
	queries *sqlc.Queries
}

func NewRepository(db sqlc.DBTX) *Repository {
	return &Repository{queries: sqlc.New(db)}
}

func (r *Repository) CreateHistory(ctx context.Context, history domainhistory.History) (*domainhistory.History, error) {
	row, err := r.queries.CreateWeatherHistory(ctx, sqlc.CreateWeatherHistoryParams{
		ID:          history.ID,
		UserID:      history.UserID,
		City:        history.City,
		Temperature: history.Temperature,
		Description: history.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("создание записи истории погоды: %w", err)
	}

	return new(toDomain(row)), nil
}

func (r *Repository) ListHistory(ctx context.Context, userID uuid.UUID, city string, limit int, offset int) ([]domainhistory.History, error) {
	rows, err := r.queries.ListWeatherHistory(ctx, sqlc.ListWeatherHistoryParams{
		UserID: userID,
		City:   city,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("получение истории погоды: %w", err)
	}

	history := make([]domainhistory.History, 0, len(rows))
	for _, row := range rows {
		history = append(history, toDomain(row))
	}

	return history, nil
}
