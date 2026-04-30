package geocoding

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/DaniilKalts/rbk-school/3-week/internal/client/geocoding/dto"
)

const (
	defaultBaseURL = "https://geocoding-api.open-meteo.com/v1/search"
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

func (c *Client) GetCoordsByCity(ctx context.Context, city string) (dto.CoordsResponse, error) {
	parsedURL, err := url.Parse(c.baseURL)
	if err != nil {
		return dto.CoordsResponse{}, fmt.Errorf("geocoding: parse base url: %w", err)
	}

	query := parsedURL.Query()
	query.Set("name", city)
	query.Set("count", "1")
	query.Set("language", "en")
	query.Set("format", "json")
	parsedURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return dto.CoordsResponse{}, fmt.Errorf("geocoding: create request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return dto.CoordsResponse{}, fmt.Errorf("geocoding: do request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return dto.CoordsResponse{}, fmt.Errorf("geocoding: unexpected status code %d", res.StatusCode)
	}

	var results dto.GeocodingResults
	if err := json.NewDecoder(res.Body).Decode(&results); err != nil {
		return dto.CoordsResponse{}, fmt.Errorf("geocoding: decode response body: %w", err)
	}

	if len(results.Results) == 0 {
		return dto.CoordsResponse{}, fmt.Errorf("geocoding: no results found for city %q", city)
	}

	return results.Results[0], nil
}
