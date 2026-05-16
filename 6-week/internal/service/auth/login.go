package auth

import (
	"context"
	"errors"

	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/user"
)

func (s *Service) Login(ctx context.Context, input LoginInput) (*Token, error) {
	email := user.NormalizeEmail(input.Email)
	if email == "" || input.Password == "" {
		return nil, ErrInvalidCredentials
	}

	u, password, err := s.repository.GetCredentialsByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	if !password.Matches(input.Password) {
		return nil, ErrInvalidCredentials
	}

	return s.generateToken(u.ID, u.Email, u.Role)
}
