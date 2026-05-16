package city

import (
	"context"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/city"
)

func (s *Service) List(ctx context.Context, userID uuid.UUID) ([]city.City, error) {
	if userID == uuid.Nil {
		return nil, city.ErrInvalidUserID
	}

	_, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	cities, err := s.repository.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return cities, nil
}
