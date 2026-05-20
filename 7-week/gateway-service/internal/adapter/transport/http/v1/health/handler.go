package health

import (
	"context"
	"net/http"

	"github.com/DaniilKalts/rbk-school/7-week/pkg/httpx"
)

type APIServiceClient interface {
	Ping(ctx context.Context) error
}

type Handler struct {
	apiServiceClient APIServiceClient
}

func NewHandler(apiServiceClient APIServiceClient) *Handler {
	return &Handler{apiServiceClient: apiServiceClient}
}

type ReadyResponse struct {
	Status     string `json:"status"`
	APIService string `json:"api_service"`
}

func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	if err := h.apiServiceClient.Ping(r.Context()); err != nil {
		httpx.JSON(w, http.StatusServiceUnavailable, ReadyResponse{
			Status:     "degraded",
			APIService: err.Error(),
		})
		return
	}

	httpx.JSON(w, http.StatusOK, ReadyResponse{
		Status:     "ok",
		APIService: "ok",
	})
}
