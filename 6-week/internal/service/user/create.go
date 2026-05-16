package user

import (
	"context"

	"github.com/DaniilKalts/rbk-school/5-week/internal/domain/user"
)

func (s *Service) Create(ctx context.Context, input CreateInput) (*user.User, error) {
	password, err := user.NewPassword(input.Password)
	if err != nil {
		return nil, err
	}

	u, err := user.NewUser(input.FirstName, input.LastName, input.Email, user.RoleUser)
	if err != nil {
		return nil, err
	}

	created, err := s.repository.Create(ctx, *u, password)
	if err != nil {
		return nil, err
	}

	return created, nil
}
