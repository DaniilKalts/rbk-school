package auth

import (
	"context"

	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/user"
)

func (s *Service) Register(ctx context.Context, input RegisterInput) (*Token, error) {
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

	return s.generateToken(created.ID, created.Email, created.Role)
}
