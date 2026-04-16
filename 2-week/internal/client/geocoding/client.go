package geocoding

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/DaniilKalts/rbk-school/2-week/internal/client/geocoding/dto"
)

const (
	defaultBaseURL = "https://geocoding-api.open-meteo.com/v1"
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

func (c *Client) GetCoordsByState(state string) (dto.CoordsResponse, error) {
	url := fmt.Sprintf("%s/search?name=%s&count=1&language=en&format=json", c.baseURL, state)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return dto.CoordsResponse{}, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return dto.CoordsResponse{}, fmt.Errorf("geocoding: failed to get coords by city: %w", err)
	}
	defer func() {
		if closeErr := res.Body.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("geocoding: failed to close response body: %w", closeErr)
		}
	}()

	if res.StatusCode != http.StatusOK {
		return dto.CoordsResponse{}, fmt.Errorf("geocoding: unexpected status code: %d", res.StatusCode)
	}

	var results dto.GeocodingResults
	if err = json.NewDecoder(res.Body).Decode(&results); err != nil {
		return dto.CoordsResponse{}, fmt.Errorf("geocoding: failed to decode response body: %w", err)
	}

	if len(results.Results) == 0 {
		return dto.CoordsResponse{}, fmt.Errorf("geocoding: no results found")
	}

	return results.Results[0], nil
}
