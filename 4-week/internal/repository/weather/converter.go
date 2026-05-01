package weather

import (
	"github.com/DaniilKalts/rbk-school/4-week/internal/adapters/database/postgres/sqlc"
	domainhistory "github.com/DaniilKalts/rbk-school/4-week/internal/domain/history"
)

func toDomain(history sqlc.WeatherHistory) domainhistory.History {
	return domainhistory.History{
		ID:          history.ID,
		UserID:      history.UserID,
		City:        history.City,
		Temperature: history.Temperature,
		Description: history.Description,
		RequestedAt: history.RequestedAt,
	}
}
