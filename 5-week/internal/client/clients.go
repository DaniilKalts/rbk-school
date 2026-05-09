package client

import (
	"net/http"

	"github.com/DaniilKalts/rbk-school/5-week/internal/client/geocoding"
	"github.com/DaniilKalts/rbk-school/5-week/internal/client/openmeteo"
)

type Clients struct {
	Geocoding *geocoding.Client
	OpenMeteo *openmeteo.Client
}

func NewClients(httpClient *http.Client) *Clients {
	return &Clients{
		Geocoding: geocoding.NewClient(httpClient),
		OpenMeteo: openmeteo.NewClient(httpClient),
	}
}
