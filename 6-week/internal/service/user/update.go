package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/user"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*user.User, error) {
	if id == uuid.Nil {
		return nil, user.ErrInvalidID
	}

	existing, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := existing.UpdateProfile(input.FirstName, input.LastName, input.Email); err != nil {
		return nil, err
	}

	updated, err := s.repository.Update(ctx, *existing)
	if err != nil {
		return nil, err
	}

	return updated, nil
}
