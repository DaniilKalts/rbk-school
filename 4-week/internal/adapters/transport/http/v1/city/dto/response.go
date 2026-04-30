package dto

import (
	"time"

	"github.com/google/uuid"
)

type CityResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	City      string    `json:"city"`
	CreatedAt time.Time `json:"created_at"`
}
