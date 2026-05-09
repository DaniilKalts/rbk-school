package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/5-week/internal/domain/user"
)

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return user.ErrInvalidID
	}

	if err := s.repository.SoftDelete(ctx, id); err != nil {
		return err
	}

	return nil
}
