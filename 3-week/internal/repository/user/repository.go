package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/DaniilKalts/rbk-school/3-week/internal/adapters/database/postgres"
	"github.com/DaniilKalts/rbk-school/3-week/internal/adapters/database/postgres/sqlc"
	domainuser "github.com/DaniilKalts/rbk-school/3-week/internal/domain/user"
)

const emailUniqueIndex = "users_email_idx"

type Repository struct {
	queries *sqlc.Queries
}

func New(db sqlc.DBTX) *Repository {
	return &Repository{queries: sqlc.New(db)}
}

func (r *Repository) Create(ctx context.Context, u domainuser.User) (*domainuser.User, error) {
	row, err := r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
	})
	if err != nil {
		if isEmailUniqueViolation(err) {
			return nil, domainuser.ErrEmailAlreadyExists
		}

		return nil, fmt.Errorf("create user: %w", err)
	}

	return new(toDomain(row)), nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*domainuser.User, error) {
	row, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainuser.ErrNotFound
		}

		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return new(toDomain(row)), nil
}

func (r *Repository) List(ctx context.Context) ([]domainuser.User, error) {
	rows, err := r.queries.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	users := make([]domainuser.User, 0, len(rows))
	for _, row := range rows {
		users = append(users, toDomain(row))
	}

	return users, nil
}

func (r *Repository) Update(ctx context.Context, u domainuser.User) (*domainuser.User, error) {
	row, err := r.queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainuser.ErrNotFound
		}

		if isEmailUniqueViolation(err) {
			return nil, domainuser.ErrEmailAlreadyExists
		}

		return nil, fmt.Errorf("update user: %w", err)
	}

	return new(toDomain(row)), nil
}

func (r *Repository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	rowsAffected, err := r.queries.SoftDeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("soft delete user: %w", err)
	}

	if rowsAffected == 0 {
		return domainuser.ErrNotFound
	}

	return nil
}

func isEmailUniqueViolation(err error) bool {
	return postgres.IsUniqueViolation(err, emailUniqueIndex)
}
