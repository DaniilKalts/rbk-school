package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/DaniilKalts/rbk-school/5-week/internal/adapters/database/postgres"
	"github.com/DaniilKalts/rbk-school/5-week/internal/adapters/database/postgres/sqlc"
	"github.com/DaniilKalts/rbk-school/5-week/internal/domain/user"
	"github.com/DaniilKalts/rbk-school/5-week/internal/service/auth"
)

const emailUniqueIndex = "users_email_idx"

type Repository struct {
	queries *sqlc.Queries
}

func NewRepository(db sqlc.DBTX) *Repository {
	return &Repository{queries: sqlc.New(db)}
}

func (r *Repository) Create(ctx context.Context, u user.User, password user.Password) (*user.User, error) {
	row, err := r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		ID:           u.ID,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Email:        u.Email,
		PasswordHash: password.Hash,
		Salt:         password.Salt,
		Role:         sqlc.UserRole(u.Role),
	})
	if err != nil {
		if isEmailUniqueViolation(err) {
			return nil, user.ErrEmailAlreadyExists
		}

		return nil, fmt.Errorf("создание пользователя: %w", err)
	}

	return new(toDomainBase(
		row.ID,
		row.FirstName,
		row.LastName,
		row.Email,
		row.Role,
		row.CreatedAt,
		row.UpdatedAt,
	)), nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	row, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrNotFound
		}

		return nil, fmt.Errorf("получение пользователя по id: %w", err)
	}

	return new(toDomainBase(
		row.ID,
		row.FirstName,
		row.LastName,
		row.Email,
		row.Role,
		row.CreatedAt,
		row.UpdatedAt,
	)), nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	row, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrNotFound
		}

		return nil, fmt.Errorf("получение пользователя по email: %w", err)
	}

	return new(toDomainBase(
		row.ID,
		row.FirstName,
		row.LastName,
		row.Email,
		row.Role,
		row.CreatedAt,
		row.UpdatedAt,
	)), nil
}

func (r *Repository) GetCredentialsByEmail(ctx context.Context, email string) (*auth.Credentials, error) {
	row, err := r.queries.GetUserCredentialsByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrNotFound
		}

		return nil, fmt.Errorf("получение учетных данных пользователя по email: %w", err)
	}

	return &auth.Credentials{
		ID:       row.ID,
		Email:    row.Email,
		Role:     user.Role(row.Role),
		Password: user.Password{Hash: row.PasswordHash, Salt: row.Salt},
	}, nil
}

func (r *Repository) List(ctx context.Context) ([]user.User, error) {
	rows, err := r.queries.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("получение списка пользователей: %w", err)
	}

	users := make([]user.User, 0, len(rows))
	for _, row := range rows {
		users = append(users, toDomainBase(
			row.ID,
			row.FirstName,
			row.LastName,
			row.Email,
			row.Role,
			row.CreatedAt,
			row.UpdatedAt,
		))
	}

	return users, nil
}

func (r *Repository) Update(ctx context.Context, u user.User) (*user.User, error) {
	row, err := r.queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrNotFound
		}

		if isEmailUniqueViolation(err) {
			return nil, user.ErrEmailAlreadyExists
		}

		return nil, fmt.Errorf("обновление пользователя: %w", err)
	}

	return new(toDomainBase(
		row.ID,
		row.FirstName,
		row.LastName,
		row.Email,
		row.Role,
		row.CreatedAt,
		row.UpdatedAt,
	)), nil
}

func (r *Repository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	rowsAffected, err := r.queries.SoftDeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("мягкое удаление пользователя: %w", err)
	}

	if rowsAffected == 0 {
		return user.ErrNotFound
	}

	return nil
}

func isEmailUniqueViolation(err error) bool {
	return postgres.IsUniqueViolation(err, emailUniqueIndex)
}
