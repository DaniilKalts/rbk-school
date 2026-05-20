package city

import (
	"github.com/DaniilKalts/rbk-school/7-week/api-service/internal/adapter/database/postgres/sqlc"
	"github.com/DaniilKalts/rbk-school/7-week/api-service/internal/domain/city"
)

func toDomain(c sqlc.UserCity) city.City {
	return city.City{
		ID:        c.ID,
		UserID:    c.UserID,
		Name:      c.City,
		CreatedAt: c.CreatedAt,
	}
}
