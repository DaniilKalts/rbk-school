package city

import (
	"context"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/city"
)

func (s *Service) Create(ctx context.Context, userID uuid.UUID, input CreateInput) (*city.City, error) {
	if userID == uuid.Nil {
		return nil, city.ErrInvalidUserID
	}

	c, err := city.NewCity(userID, input.Name)
	if err != nil {
		return nil, err
	}

	_, err = s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	created, err := s.repository.Create(ctx, *c)
	if err != nil {
		return nil, err
	}

	return created, nil
}
