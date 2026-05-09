package auth

import (
	"context"
	"errors"

	"github.com/DaniilKalts/rbk-school/5-week/internal/domain/user"
)

func (s *Service) Login(ctx context.Context, input LoginInput) (*Token, error) {
	email := user.NormalizeEmail(input.Email)
	if email == "" || input.Password == "" {
		return nil, ErrInvalidCredentials
	}

	credentials, err := s.repository.GetCredentialsByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	if !credentials.Verify(input.Password) {
		return nil, ErrInvalidCredentials
	}

	return s.generateToken(credentials.ID, credentials.Email, credentials.Role)
}
