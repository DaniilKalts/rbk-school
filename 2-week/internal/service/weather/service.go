package weather

import (
	"sort"
	"strconv"
	"strings"

	cscDTO "github.com/DaniilKalts/rbk-school/2-week/internal/client/countrystatecity/dto"
	geoDTO "github.com/DaniilKalts/rbk-school/2-week/internal/client/geocoding/dto"
	weatherDTO "github.com/DaniilKalts/rbk-school/2-week/internal/client/openmeteo/dto"
)

type CityListClient interface {
	GetStatesByCountry(countryCode string) ([]cscDTO.StateResponse, error)
}

type GeocodingClient interface {
	GetCoordsByState(countryCode string) (geoDTO.CoordsResponse, error)
}

type WeatherClient interface {
	GetWeatherByCoords(latitude, longitude float64) (weatherDTO.WeatherResponse, error)
}

type Service struct {
	cityListClient  CityListClient
	geocodingClient GeocodingClient
	weatherClient   WeatherClient
}

func NewService(cityListClient CityListClient, geocodingClient GeocodingClient, weatherClient WeatherClient) *Service {
	return &Service{
		cityListClient:  cityListClient,
		geocodingClient: geocodingClient,
		weatherClient:   weatherClient,
	}
}

func (s *Service) GetWeatherByCity(city string) (weatherDTO.WeatherResponse, error) {
	coords, err := s.geocodingClient.GetCoordsByState(city)
	if err != nil {
		return weatherDTO.WeatherResponse{}, err
	}

	weather, err := s.weatherClient.GetWeatherByCoords(coords.Latitude, coords.Longitude)
	if err != nil {
		return weatherDTO.WeatherResponse{}, err
	}
	weather.City = normalizeCityName(city)

	return weather, nil
}

func (s *Service) GetWeatherByCountry(countryCode string) ([]weatherDTO.WeatherResponse, error) {
	states, err := s.cityListClient.GetStatesByCountry(countryCode)
	if err != nil {
		return []weatherDTO.WeatherResponse{}, err
	}

	weathers := make([]weatherDTO.WeatherResponse, 0, len(states))
	for _, state := range states {
		lat, err := strconv.ParseFloat(state.Latitude, 64)
		if err != nil {
			return []weatherDTO.WeatherResponse{}, err
		}
		lon, err := strconv.ParseFloat(state.Longitude, 64)
		if err != nil {
			return []weatherDTO.WeatherResponse{}, err
		}

		weather, err := s.weatherClient.GetWeatherByCoords(lat, lon)
		if err != nil {
			continue
		}
		weather.City = normalizeCityName(state.Name)

		weathers = append(weathers, weather)
	}

	return weathers, nil
}

func (s *Service) GetTopWarmestCities(countryCode string, limit int) ([]weatherDTO.WeatherResponse, error) {
	weathers, err := s.GetWeatherByCountry(countryCode)
	if err != nil {
		return nil, err
	}

	sort.Slice(weathers, func(i, j int) bool {
		return weathers[i].Current.Temperature2M > weathers[j].Current.Temperature2M
	})

	if len(weathers) < limit {
		limit = len(weathers)
	}

	return weathers[:limit], nil
}

func normalizeCityName(city string) string {
	city = strings.TrimSpace(city)
	if city == "" {
		return city
	}

	city = strings.ToLower(city)
	return strings.ToUpper(city[:1]) + city[1:]
}
