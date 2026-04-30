package user

import (
	"time"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/3-week/internal/adapters/database/postgres/sqlc"
	domainuser "github.com/DaniilKalts/rbk-school/3-week/internal/domain/user"
)

func toDomainBase(
	id uuid.UUID,
	firstName string,
	lastName string,
	email string,
	role sqlc.UserRole,
	createdAt time.Time,
	updatedAt time.Time,
) domainuser.User {
	return domainuser.User{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Role:      domainuser.Role(role),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
