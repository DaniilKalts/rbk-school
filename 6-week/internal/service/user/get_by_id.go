package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/5-week/internal/domain/user"
)

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	if id == uuid.Nil {
		return nil, user.ErrInvalidID
	}

	u, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return u, nil
}
