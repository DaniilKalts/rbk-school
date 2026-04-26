package dto

import (
	"github.com/google/uuid"

	domainweather "github.com/DaniilKalts/rbk-school/3-week/internal/domain/weather"
)

func ToUserWeatherResponse(userID uuid.UUID, weathers []domainweather.Weather) UserWeatherResponse {
	responses := make([]WeatherResponse, 0, len(weathers))
	for _, weather := range weathers {
		responses = append(responses, ToWeatherResponse(weather))
	}

	return UserWeatherResponse{
		UserID:  userID,
		Weather: responses,
	}
}

func ToWeatherResponse(weather domainweather.Weather) WeatherResponse {
	return WeatherResponse{
		City:        weather.City,
		Temperature: weather.Temperature,
		FeelsLike:   weather.FeelsLike,
		Description: weather.Description,
		RequestedAt: weather.RequestedAt,
	}
}

func ToUserWeatherHistoryResponse(userID uuid.UUID, city string, history []domainweather.History) UserWeatherHistoryResponse {
	responses := make([]WeatherHistoryResponse, 0, len(history))
	includeCity := city == ""
	for _, item := range history {
		responses = append(responses, ToWeatherHistoryResponse(item, includeCity))
	}

	return UserWeatherHistoryResponse{
		UserID:  userID,
		City:    city,
		History: responses,
	}
}

func ToWeatherHistoryResponse(history domainweather.History, includeCity bool) WeatherHistoryResponse {
	response := WeatherHistoryResponse{
		Temperature: history.Temperature,
		Description: history.Description,
		RequestedAt: history.RequestedAt,
	}
	if includeCity {
		response.City = history.City
	}

	return response
}
