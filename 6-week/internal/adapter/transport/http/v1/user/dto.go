package user

import (
	"time"

	"github.com/DaniilKalts/rbk-school/5-week/internal/domain/user"

	serviceuser "github.com/DaniilKalts/rbk-school/5-week/internal/service/user"
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

func ToCreateInput(r CreateUserRequest) serviceuser.CreateInput {
	return serviceuser.CreateInput{FirstName: r.FirstName, LastName: r.LastName, Email: r.Email, Password: r.Password}
}

func ToUpdateInput(r UpdateUserRequest) serviceuser.UpdateInput {
	return serviceuser.UpdateInput{FirstName: r.FirstName, LastName: r.LastName, Email: r.Email}
}

func ToUserResponse(u user.User) UserResponse {
	return UserResponse{ID: u.ID.String(), FirstName: u.FirstName, LastName: u.LastName, Email: u.Email, CreatedAt: u.CreatedAt, UpdatedAt: u.UpdatedAt}
}

func ToUserResponses(users []user.User) []UserResponse {
	res := make([]UserResponse, 0, len(users))
	for _, u := range users {
		res = append(res, ToUserResponse(u))
	}
	return res
}
