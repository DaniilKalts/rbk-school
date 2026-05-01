package city

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/4-week/internal/adapters/database/postgres"
	"github.com/DaniilKalts/rbk-school/4-week/internal/adapters/database/postgres/sqlc"
	domaincity "github.com/DaniilKalts/rbk-school/4-week/internal/domain/city"
)

const userCityUniqueConstraint = "user_cities_user_id_city_key"

type Repository struct {
	queries *sqlc.Queries
}

func New(db sqlc.DBTX) *Repository {
	return &Repository{queries: sqlc.New(db)}
}

func (r *Repository) Create(ctx context.Context, c domaincity.City) (*domaincity.City, error) {
	row, err := r.queries.CreateUserCity(ctx, sqlc.CreateUserCityParams{
		ID:     c.ID,
		UserID: c.UserID,
		City:   c.Name,
	})
	if err != nil {
		if postgres.IsUniqueViolation(err, userCityUniqueConstraint) {
			return nil, domaincity.ErrAlreadyExists
		}

		return nil, fmt.Errorf("create user city: %w", err)
	}

	return new(toDomain(row)), nil
}

func (r *Repository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]domaincity.City, error) {
	rows, err := r.queries.ListUserCities(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list user cities: %w", err)
	}

	cities := make([]domaincity.City, 0, len(rows))
	for _, row := range rows {
		cities = append(cities, toDomain(row))
	}

	return cities, nil
}

func (r *Repository) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	rowsAffected, err := r.queries.DeleteUserCity(ctx, sqlc.DeleteUserCityParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		return fmt.Errorf("delete user city: %w", err)
	}

	if rowsAffected == 0 {
		return domaincity.ErrNotFound
	}

	return nil
}
