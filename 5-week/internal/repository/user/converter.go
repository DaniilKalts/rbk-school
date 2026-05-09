package user

import (
	"time"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/5-week/internal/adapter/database/postgres/sqlc"
	"github.com/DaniilKalts/rbk-school/5-week/internal/domain/user"
)

func toDomainBase(
	id uuid.UUID,
	firstName string,
	lastName string,
	email string,
	role sqlc.UserRole,
	createdAt time.Time,
	updatedAt time.Time,
) user.User {
	return user.User{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Role:      user.Role(role),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
