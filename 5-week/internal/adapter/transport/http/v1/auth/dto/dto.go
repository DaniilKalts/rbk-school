package dto

import (
	"github.com/DaniilKalts/rbk-school/5-week/internal/service/auth"
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

func ToRegisterInput(req RegisterRequest) auth.RegisterInput {
	return auth.RegisterInput{FirstName: req.FirstName, LastName: req.LastName, Email: req.Email, Password: req.Password}
}

func ToLoginInput(req LoginRequest) auth.LoginInput {
	return auth.LoginInput{Email: req.Email, Password: req.Password}
}

func ToTokenResponse(token auth.Token) TokenResponse {
	return TokenResponse{AccessToken: token.AccessToken, ExpiresAt: token.ExpiresAt.UTC().Format("2006-01-02T15:04:05Z")}
}
