package auth

import "context"

func (s *Service) Logout(ctx context.Context, accessToken string) error {
	return s.tokenManager.Revoke(ctx, accessToken)
}
