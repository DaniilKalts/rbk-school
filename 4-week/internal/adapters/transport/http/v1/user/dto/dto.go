package dto

import (
	"time"

	domainuser "github.com/DaniilKalts/rbk-school/4-week/internal/domain/user"
	serviceuser "github.com/DaniilKalts/rbk-school/4-week/internal/service/user"
)

type CreateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}
type UpdateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
type UserResponse struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToCreateInput(req CreateUserRequest) serviceuser.CreateInput {
	return serviceuser.CreateInput{FirstName: req.FirstName, LastName: req.LastName, Email: req.Email, Password: req.Password}
}
func ToUpdateInput(req UpdateUserRequest) serviceuser.UpdateInput {
	return serviceuser.UpdateInput{FirstName: req.FirstName, LastName: req.LastName, Email: req.Email}
}
func ToUserResponse(u domainuser.User) UserResponse {
	return UserResponse{ID: u.ID.String(), FirstName: u.FirstName, LastName: u.LastName, Email: u.Email, CreatedAt: u.CreatedAt, UpdatedAt: u.UpdatedAt}
}
func ToUserResponses(users []domainuser.User) []UserResponse {
	res := make([]UserResponse, 0, len(users))
	for _, u := range users {
		res = append(res, ToUserResponse(u))
	}
	return res
}
