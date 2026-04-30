package dto

import (
	domaincity "github.com/DaniilKalts/rbk-school/3-week/internal/domain/city"
	servicecity "github.com/DaniilKalts/rbk-school/3-week/internal/service/city"
)

func ToCreateInput(req CreateCityRequest) servicecity.CreateInput {
	return servicecity.CreateInput{Name: req.City}
}

func ToCityResponse(c domaincity.City) CityResponse {
	return CityResponse{
		ID:        c.ID,
		UserID:    c.UserID,
		City:      c.Name,
		CreatedAt: c.CreatedAt,
	}
}

func ToCityResponses(cities []domaincity.City) []CityResponse {
	responses := make([]CityResponse, 0, len(cities))
	for _, c := range cities {
		responses = append(responses, ToCityResponse(c))
	}

	return responses
}
