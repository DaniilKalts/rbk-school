package city

import (
	"github.com/DaniilKalts/rbk-school/3-week/internal/adapters/database/postgres/sqlc"
	domaincity "github.com/DaniilKalts/rbk-school/3-week/internal/domain/city"
)

func toDomain(c sqlc.UserCity) domaincity.City {
	return domaincity.City{
		ID:        c.ID,
		UserID:    c.UserID,
		Name:      c.City,
		CreatedAt: c.CreatedAt,
	}
}
