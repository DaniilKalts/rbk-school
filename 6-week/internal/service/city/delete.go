package city

import (
	"context"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/city"
)

func (s *Service) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	if id == uuid.Nil {
		return city.ErrInvalidID
	}

	if userID == uuid.Nil {
		return city.ErrInvalidUserID
	}

	_, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	return s.repository.Delete(ctx, userID, id)
}
