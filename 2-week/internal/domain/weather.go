package domain

import (
	"errors"
	"fmt"
	"strings"
)

var ErrEmptyCityName = errors.New("city name is required")

type Weather struct {
	City       string
	Conditions Conditions
}

type Conditions struct {
	Temperature float64
	FeelsLike   float64
}

func NewWeather(city string, latitude, longitude float64, conditions Conditions) (Weather, error) {
	city = normalizeCityName(city)
	if city == "" {
		return Weather{}, ErrEmptyCityName
	}

	if err := validateCoordinates(latitude, longitude); err != nil {
		return Weather{}, err
	}

	return Weather{
		City: city,
		Conditions: Conditions{
			Temperature: conditions.Temperature,
			FeelsLike:   conditions.FeelsLike,
		},
	}, nil
}

func (w Weather) Recommendation() string {
	switch t := w.Conditions.Temperature; {
	case t < 10:
		return "It is cold outside, wear warm clothes."
	case t < 20:
		return "It is cool outside, a jacket will help."
	default:
		return "It is warm outside, light clothes are enough."
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
