package city

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/7-week/internal/adapter/database/postgres"
	"github.com/DaniilKalts/rbk-school/7-week/internal/adapter/database/postgres/sqlc"
	"github.com/DaniilKalts/rbk-school/7-week/internal/domain/city"
)

const userCityUniqueConstraint = "user_cities_user_id_city_key"

type Repository struct {
	queries *sqlc.Queries
}

func NewRepository(db sqlc.DBTX) *Repository {
	return &Repository{queries: sqlc.New(db)}
}

func (r *Repository) Create(ctx context.Context, c city.City) (*city.City, error) {
	row, err := r.queries.CreateUserCity(ctx, sqlc.CreateUserCityParams{
		ID:     c.ID,
		UserID: c.UserID,
		City:   c.Name,
	})
	if err != nil {
		if postgres.IsUniqueViolation(err, userCityUniqueConstraint) {
			return nil, city.ErrAlreadyExists
		}

		return nil, fmt.Errorf("создание города пользователя: %w", err)
	}

	return new(toDomain(row)), nil
}

func (r *Repository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]city.City, error) {
	rows, err := r.queries.ListUserCities(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("получение списка городов пользователя: %w", err)
	}

	cities := make([]city.City, 0, len(rows))
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
		return fmt.Errorf("удаление города пользователя: %w", err)
	}

	if rowsAffected == 0 {
		return city.ErrNotFound
	}

	return nil
}
