package auth

import (
	"time"

	"github.com/DaniilKalts/rbk-school/7-week/api-service/internal/service/auth"
)

type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   string `json:"expires_at"`
}

func ToRegisterInput(r RegisterRequest) auth.RegisterInput {
	return auth.RegisterInput{FirstName: r.FirstName, LastName: r.LastName, Email: r.Email, Password: r.Password}
}

func ToLoginInput(r LoginRequest) auth.LoginInput {
	return auth.LoginInput{Email: r.Email, Password: r.Password}
}

func ToTokenResponse(token auth.Token) TokenResponse {
	return TokenResponse{AccessToken: token.AccessToken, ExpiresAt: token.ExpiresAt.UTC().Format(time.RFC3339)}
}
