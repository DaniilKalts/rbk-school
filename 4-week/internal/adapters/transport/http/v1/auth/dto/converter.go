package dto

import serviceauth "github.com/DaniilKalts/rbk-school/3-week/internal/service/auth"

func ToRegisterInput(req RegisterRequest) serviceauth.RegisterInput {
	return serviceauth.RegisterInput{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
	}
}

func ToLoginInput(req LoginRequest) serviceauth.LoginInput {
	return serviceauth.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}
}

func ToTokenResponse(token serviceauth.Token) TokenResponse {
	return TokenResponse{
		AccessToken: token.AccessToken,
		ExpiresAt:   token.ExpiresAt,
	}
}
