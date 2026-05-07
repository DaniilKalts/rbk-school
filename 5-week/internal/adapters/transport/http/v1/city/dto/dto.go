package dto

import (
	"time"

	domaincity "github.com/DaniilKalts/rbk-school/5-week/internal/domain/city"
	servicecity "github.com/DaniilKalts/rbk-school/5-week/internal/service/city"
)

type CreateCityRequest struct {
	City string `json:"city"`
}

type CityResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	City      string    `json:"city"`
	CreatedAt time.Time `json:"created_at"`
}

func ToCreateInput(req CreateCityRequest) servicecity.CreateInput {
	return servicecity.CreateInput{Name: req.City}
}
func ToCityResponse(c domaincity.City) CityResponse {
	return CityResponse{ID: c.ID.String(), UserID: c.UserID.String(), City: c.Name, CreatedAt: c.CreatedAt}
}
func ToCityResponses(cities []domaincity.City) []CityResponse {
	res := make([]CityResponse, 0, len(cities))
	for _, c := range cities {
		res = append(res, ToCityResponse(c))
	}
	return res
}
