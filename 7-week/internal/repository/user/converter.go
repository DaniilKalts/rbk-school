package user

import (
	"time"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/7-week/internal/adapter/database/postgres/sqlc"
	"github.com/DaniilKalts/rbk-school/7-week/internal/domain/user"
)

type userRow struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Email     string
	Role      sqlc.UserRole
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func toDomain(r userRow) *user.User {
	return &user.User{
		ID:        r.ID,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Email:     r.Email,
		Role:      user.Role(r.Role),
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
