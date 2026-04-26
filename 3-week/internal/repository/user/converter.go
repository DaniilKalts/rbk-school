package user

import (
	"github.com/DaniilKalts/rbk-school/3-week/internal/adapters/database/postgres/sqlc"
	domainuser "github.com/DaniilKalts/rbk-school/3-week/internal/domain/user"
)

func toDomain(u sqlc.User) domainuser.User {
	return domainuser.User{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
