package history

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type History struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	City        string
	Temperature float64
	Description string
	RequestedAt time.Time
}

func New(id uuid.UUID, userID uuid.UUID, city string, temperature float64, description string) (*History, error) {
	if id == uuid.Nil {
		return nil, ErrInvalidID
	}
	if userID == uuid.Nil {
		return nil, ErrInvalidUserID
	}
	
	city = strings.TrimSpace(city)
	if city == "" {
		return nil, ErrInvalidCity
	}

	return &History{
		ID:          id,
		UserID:      userID,
		City:        city,
		Temperature: temperature,
		Description: strings.TrimSpace(description),
	}, nil
}
