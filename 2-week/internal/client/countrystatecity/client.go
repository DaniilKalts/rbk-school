package countrystatecity

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/DaniilKalts/rbk-school/2-week/internal/client/countrystatecity/dto"
)

const (
	defaultBaseURL = "https://api.countrystatecity.in/v1"
	defaultTimeout = 15 * time.Second
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

func NewClient(httpClient *http.Client, apiKey string) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}

	return &Client{
		httpClient: httpClient,
		baseURL:    defaultBaseURL,
		apiKey:     apiKey,
	}
}

func (c *Client) GetStatesByCountry(ctx context.Context, countryCode string) ([]dto.StateResponse, error) {
	url := fmt.Sprintf("%s/countries/%s/states", c.baseURL, countryCode)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("countrystatecity: create request: %w", err)
	}
	req.Header.Set("X-CSCAPI-KEY", c.apiKey)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("countrystatecity: do request: %w", err)
	}
	defer func() {
		if closeErr := res.Body.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("countrystatecity: close response body: %w", closeErr)
		}
	}()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("countrystatecity: unexpected status code %d", res.StatusCode)
	}

	var states []dto.StateResponse
	if err = json.NewDecoder(res.Body).Decode(&states); err != nil {
		return nil, fmt.Errorf("countrystatecity: decode response body: %w", err)
	}

	return states, nil
}
