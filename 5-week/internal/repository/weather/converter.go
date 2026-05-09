package weather

import (
	"github.com/DaniilKalts/rbk-school/5-week/internal/adapters/database/postgres/sqlc"
	"github.com/DaniilKalts/rbk-school/5-week/internal/domain/history"
)

func toDomain(row sqlc.WeatherHistory) history.History {
	return history.History{
		ID:          row.ID,
		UserID:      row.UserID,
		City:        row.City,
		Temperature: row.Temperature,
		Description: row.Description,
		RequestedAt: row.RequestedAt,
	}
}
