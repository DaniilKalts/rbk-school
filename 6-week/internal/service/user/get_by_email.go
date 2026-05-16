package user

import (
	"context"

	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/user"
)

func (s *Service) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	email = user.NormalizeEmail(email)
	if email == "" {
		return nil, user.ErrInvalidEmail
	}

	u, err := s.repository.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return u, nil
}
