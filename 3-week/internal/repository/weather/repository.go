package weather

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/3-week/internal/adapters/database/postgres/sqlc"
	domainweather "github.com/DaniilKalts/rbk-school/3-week/internal/domain/weather"
)

type Repository struct {
	queries *sqlc.Queries
}

func New(db sqlc.DBTX) *Repository {
	return &Repository{queries: sqlc.New(db)}
}

func (r *Repository) CreateHistory(ctx context.Context, history domainweather.History) (*domainweather.History, error) {
	row, err := r.queries.CreateWeatherHistory(ctx, sqlc.CreateWeatherHistoryParams{
		ID:          history.ID,
		UserID:      history.UserID,
		City:        history.City,
		Temperature: history.Temperature,
		Description: history.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("create weather history: %w", err)
	}

	return new(toDomain(row)), nil
}

func (r *Repository) ListHistoryByUserAndCity(ctx context.Context, userID uuid.UUID, city string, limit int) ([]domainweather.History, error) {
	rows, err := r.queries.ListWeatherHistoryByUserAndCity(ctx, sqlc.ListWeatherHistoryByUserAndCityParams{
		UserID: userID,
		City:   city,
		Limit:  int32(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("list weather history by user and city: %w", err)
	}

	history := make([]domainweather.History, 0, len(rows))
	for _, row := range rows {
		history = append(history, toDomain(row))
	}

	return history, nil
}

func toDomain(history sqlc.WeatherHistory) domainweather.History {
	return domainweather.History{
		ID:          history.ID,
		UserID:      history.UserID,
		City:        history.City,
		Temperature: history.Temperature,
		Description: history.Description,
		RequestedAt: history.RequestedAt,
	}
}
