package dto

import (
	"time"

	"github.com/google/uuid"
)

type UserWeatherResponse struct {
	UserID  uuid.UUID         `json:"user_id"`
	Weather []WeatherResponse `json:"weather"`
}

type WeatherResponse struct {
	City        string    `json:"city"`
	Temperature float64   `json:"temperature"`
	FeelsLike   float64   `json:"feels_like"`
	Description string    `json:"description"`
	RequestedAt time.Time `json:"requested_at"`
}

type UserWeatherHistoryResponse struct {
	UserID  uuid.UUID                `json:"user_id"`
	City    string                   `json:"city"`
	History []WeatherHistoryResponse `json:"history"`
}

type WeatherHistoryResponse struct {
	Temperature float64   `json:"temperature"`
	Description string    `json:"description"`
	RequestedAt time.Time `json:"requested_at"`
}
