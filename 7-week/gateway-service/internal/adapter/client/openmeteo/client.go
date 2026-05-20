package openmeteo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/DaniilKalts/rbk-school/7-week/gateway-service/internal/adapter/client/openmeteo/dto"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient(httpClient *http.Client, baseURL string) *Client {
	return &Client{
		httpClient: httpClient,
		baseURL:    baseURL,
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
