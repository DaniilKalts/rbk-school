package dto

import (
	domainuser "github.com/DaniilKalts/rbk-school/3-week/internal/domain/user"
	serviceuser "github.com/DaniilKalts/rbk-school/3-week/internal/service/user"
)

func ToCreateInput(req CreateUserRequest) serviceuser.CreateInput {
	return serviceuser.CreateInput{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
	}
}

func ToUpdateInput(req UpdateUserRequest) serviceuser.UpdateInput {
	return serviceuser.UpdateInput{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}
}

func ToUserResponse(u domainuser.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func ToUserResponses(users []domainuser.User) []UserResponse {
	responses := make([]UserResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, ToUserResponse(u))
	}

	return responses
}
