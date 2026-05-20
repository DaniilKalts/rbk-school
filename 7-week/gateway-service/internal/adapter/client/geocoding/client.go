package geocoding

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/DaniilKalts/rbk-school/7-week/gateway-service/internal/adapter/client/geocoding/dto"
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

func (c *Client) GetCoordsByCity(ctx context.Context, city string) (dto.CoordsResponse, error) {
	parsedURL, err := url.Parse(c.baseURL)
	if err != nil {
		return dto.CoordsResponse{}, fmt.Errorf("geocoding: ошибка разбора базового URL: %w", err)
	}

	query := parsedURL.Query()
	query.Set("name", city)
	query.Set("count", "1")
	query.Set("language", "en")
	query.Set("format", "json")
	parsedURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return dto.CoordsResponse{}, fmt.Errorf("geocoding: ошибка создания запроса: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return dto.CoordsResponse{}, fmt.Errorf("geocoding: ошибка выполнения запроса: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return dto.CoordsResponse{}, fmt.Errorf("geocoding: неожиданный статус-код %d", res.StatusCode)
	}

	var results dto.GeocodingResults
	if err := json.NewDecoder(res.Body).Decode(&results); err != nil {
		return dto.CoordsResponse{}, fmt.Errorf("geocoding: ошибка декодирования тела ответа: %w", err)
	}

	if len(results.Results) == 0 {
		return dto.CoordsResponse{}, fmt.Errorf("geocoding: не найдено результатов для города %q", city)
	}

	return results.Results[0], nil
}
