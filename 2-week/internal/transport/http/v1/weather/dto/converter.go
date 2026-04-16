package dto

import "github.com/DaniilKalts/rbk-school/2-week/internal/domain"

func FromDomainWeather(weather domain.Weather) WeatherResponse {
	return WeatherResponse{
		City:           weather.City,
		Recommendation: weather.Recommendation(),
		Conditions: Conditions{
			Temperature: weather.Conditions.Temperature,
			FeelsLike:   weather.Conditions.FeelsLike,
		},
	}
}

func FromDomainWeathers(weathers []domain.Weather) []WeatherResponse {
	responses := make([]WeatherResponse, 0, len(weathers))
	for _, weather := range weathers {
		responses = append(responses, FromDomainWeather(weather))
	}

	return responses
}
