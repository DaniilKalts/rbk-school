package user

import (
	"context"

	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/user"
)

func (s *Service) List(ctx context.Context) ([]user.User, error) {
	users, err := s.repository.List(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}
