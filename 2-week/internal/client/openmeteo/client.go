package openmeteo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/DaniilKalts/rbk-school/2-week/internal/client/openmeteo/dto"
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

func (c *Client) GetWeatherByCoords(latitude, longitude float64) (dto.WeatherResponse, error) {
	parsedURL, err := url.Parse(c.baseURL)
	if err != nil {
		return dto.WeatherResponse{}, fmt.Errorf("openmeteo: parse base url: %w", err)
	}

	queryParams := parsedURL.Query()
	queryParams.Set("latitude", strconv.FormatFloat(latitude, 'f', -1, 64))
	queryParams.Set("longitude", strconv.FormatFloat(longitude, 'f', -1, 64))
	queryParams.Set("current", "temperature_2m,apparent_temperature")

	parsedURL.RawQuery = queryParams.Encode()

	req, err := http.NewRequest(http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return dto.WeatherResponse{}, fmt.Errorf("openmeteo: create request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return dto.WeatherResponse{}, fmt.Errorf("openmeteo: do request: %w", err)
	}
	defer func() {
		if closeErr := res.Body.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("openmeteo: close response body: %w", closeErr)
		}
	}()

	if res.StatusCode != http.StatusOK {
		return dto.WeatherResponse{}, fmt.Errorf("openmeteo: unexpected status code %d", res.StatusCode)
	}

	var weather dto.WeatherResponse
	if err := json.NewDecoder(res.Body).Decode(&weather); err != nil {
		return dto.WeatherResponse{}, fmt.Errorf("openmeteo: decode response body: %w", err)
	}

	return weather, nil
}
