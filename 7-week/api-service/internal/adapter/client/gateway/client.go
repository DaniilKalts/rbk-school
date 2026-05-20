package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/DaniilKalts/rbk-school/7-week/api-service/internal/adapter/client/gateway/dto"
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

func (c *Client) GetWeatherByCity(ctx context.Context, city string) (dto.WeatherResponse, error) {
	parsedURL, err := url.Parse(c.baseURL + "/api/v1/weather")
	if err != nil {
		return dto.WeatherResponse{}, fmt.Errorf("gateway: ошибка разбора URL: %w", err)
	}

	query := parsedURL.Query()
	query.Set("city", city)
	parsedURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return dto.WeatherResponse{}, fmt.Errorf("gateway: ошибка создания запроса: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return dto.WeatherResponse{}, fmt.Errorf("gateway: ошибка выполнения запроса: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return dto.WeatherResponse{}, fmt.Errorf("gateway: неожиданный статус-код %d", res.StatusCode)
	}

	var weather dto.WeatherResponse
	if err := json.NewDecoder(res.Body).Decode(&weather); err != nil {
		return dto.WeatherResponse{}, fmt.Errorf("gateway: ошибка декодирования тела ответа: %w", err)
	}

	return weather, nil
}
