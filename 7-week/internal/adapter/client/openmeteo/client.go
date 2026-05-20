package openmeteo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/DaniilKalts/rbk-school/7-week/internal/adapter/client/openmeteo/dto"
)

const (
	defaultBaseURL = "https://api.open-meteo.com/v1/forecast"
	defaultTimeout = 15 * time.Second
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}

	return &Client{
		httpClient: httpClient,
		baseURL:    defaultBaseURL,
	}
}

func (c *Client) GetWeatherByCoords(ctx context.Context, latitude, longitude float64) (dto.WeatherResponse, error) {
	parsedURL, err := url.Parse(c.baseURL)
	if err != nil {
		return dto.WeatherResponse{}, fmt.Errorf("openmeteo: ошибка разбора базового URL: %w", err)
	}

	query := parsedURL.Query()
	query.Set("latitude", strconv.FormatFloat(latitude, 'f', -1, 64))
	query.Set("longitude", strconv.FormatFloat(longitude, 'f', -1, 64))
	query.Set("current", "temperature_2m,apparent_temperature,weather_code")
	parsedURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return dto.WeatherResponse{}, fmt.Errorf("openmeteo: ошибка создания запроса: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return dto.WeatherResponse{}, fmt.Errorf("openmeteo: ошибка выполнения запроса: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return dto.WeatherResponse{}, fmt.Errorf("openmeteo: неожиданный статус-код %d", res.StatusCode)
	}

	var weather dto.WeatherResponse
	if err := json.NewDecoder(res.Body).Decode(&weather); err != nil {
		return dto.WeatherResponse{}, fmt.Errorf("openmeteo: ошибка декодирования тела ответа: %w", err)
	}

	return weather, nil
}
