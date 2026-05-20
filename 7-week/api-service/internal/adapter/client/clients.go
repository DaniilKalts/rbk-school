package client

import (
	"net/http"

	"github.com/DaniilKalts/rbk-school/7-week/api-service/internal/adapter/client/gateway"
)

type Clients struct {
	Gateway *gateway.Client
}

func NewClients(httpClient *http.Client, gatewayBaseURL string) *Clients {
	return &Clients{
		Gateway: gateway.NewClient(httpClient, gatewayBaseURL),
	}
}
