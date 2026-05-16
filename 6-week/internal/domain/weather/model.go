package weather

import (
	"time"

	domaincity "github.com/DaniilKalts/rbk-school/6-week/internal/domain/city"
)

const (
	temperatureMinCelsius = -100.0
	temperatureMaxCelsius = 60.0
)

type Weather struct {
	City        string
	Temperature float64
	FeelsLike   float64
	Description string
	RequestedAt time.Time
}

func NewWeather(city string, temperature, feelsLike float64, weatherCode int) (Weather, error) {
	city = domaincity.NormalizeCityName(city)
	if city == "" {
		return Weather{}, ErrInvalidCity
	}
	if !validTemperature(temperature) || !validTemperature(feelsLike) {
		return Weather{}, ErrInvalidTemperature
	}

	return Weather{
		City:        city,
		Temperature: temperature,
		FeelsLike:   feelsLike,
		Description: descriptionByCode(weatherCode),
	}, nil
}

func validTemperature(t float64) bool {
	return t >= temperatureMinCelsius && t <= temperatureMaxCelsius
}

func descriptionByCode(code int) string {
	switch code {
	case 0:
		return "clear sky"
	case 1, 2, 3:
		return "partly cloudy"
	case 45, 48:
		return "fog"
	case 51, 53, 55, 56, 57:
		return "drizzle"
	case 61, 63, 65, 66, 67:
		return "rain"
	case 71, 73, 75, 77:
		return "snow"
	case 80, 81, 82:
		return "rain showers"
	case 85, 86:
		return "snow showers"
	case 95, 96, 99:
		return "thunderstorm"
	default:
		return "unknown"
	}
}
