package weather

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Weather struct {
	City        string
	Temperature float64
	FeelsLike   float64
	Description string
	RequestedAt time.Time
}

type History struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	City        string
	Temperature float64
	Description string
	RequestedAt time.Time
}

func New(city string, latitude, longitude float64, temperature, feelsLike float64, weatherCode int) (Weather, error) {
	city = normalizeCityName(city)
	if city == "" {
		return Weather{}, ErrInvalidCity
	}

	if err := validateCoordinates(latitude, longitude); err != nil {
		return Weather{}, err
	}

	return Weather{
		City:        city,
		Temperature: temperature,
		FeelsLike:   feelsLike,
		Description: DescriptionByCode(weatherCode),
	}, nil
}

func DescriptionByCode(code int) string {
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

func validateCoordinates(latitude, longitude float64) error {
	if latitude < -90 || latitude > 90 {
		return fmt.Errorf("invalid latitude: %.6f", latitude)
	}
	if longitude < -180 || longitude > 180 {
		return fmt.Errorf("invalid longitude: %.6f", longitude)
	}

	return nil
}

func normalizeCityName(city string) string {
	city = strings.TrimSpace(city)
	if city == "" {
		return city
	}
	city = strings.ToLower(city)

	return strings.ToUpper(city[:1]) + city[1:]
}
