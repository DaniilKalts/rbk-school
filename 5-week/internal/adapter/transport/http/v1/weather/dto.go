package weather

import (
	"time"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/5-week/internal/domain/history"
	"github.com/DaniilKalts/rbk-school/5-week/internal/domain/weather"
)

type WeatherResponse struct {
	City        string    `json:"city"`
	Temperature float64   `json:"temperature"`
	FeelsLike   float64   `json:"feels_like"`
	Description string    `json:"description"`
	RequestedAt time.Time `json:"requested_at"`
}
type UserWeatherResponse struct {
	UserID  string            `json:"user_id"`
	Weather []WeatherResponse `json:"weather"`
}
type WeatherHistoryResponse struct {
	City        string    `json:"city,omitempty"`
	Temperature float64   `json:"temperature"`
	Description string    `json:"description"`
	RequestedAt time.Time `json:"requested_at"`
}
type UserWeatherHistoryResponse struct {
	UserID  string                   `json:"user_id"`
	City    string                   `json:"city,omitempty"`
	History []WeatherHistoryResponse `json:"history"`
}

func ToUserWeatherResponse(userID uuid.UUID, weather []weather.Weather) UserWeatherResponse {
	items := make([]WeatherResponse, 0, len(weather))
	for _, w := range weather {
		items = append(items, WeatherResponse{City: w.City, Temperature: w.Temperature, FeelsLike: w.FeelsLike, Description: w.Description, RequestedAt: w.RequestedAt})
	}
	return UserWeatherResponse{UserID: userID.String(), Weather: items}
}

func ToUserWeatherHistoryResponse(userID uuid.UUID, city string, history []history.History) UserWeatherHistoryResponse {
	items := make([]WeatherHistoryResponse, 0, len(history))
	for _, h := range history {
		items = append(items, WeatherHistoryResponse{City: h.City, Temperature: h.Temperature, Description: h.Description, RequestedAt: h.RequestedAt})
	}
	res := UserWeatherHistoryResponse{UserID: userID.String(), History: items}
	if city != "" {
		res.City = city
	}
	return res
}
