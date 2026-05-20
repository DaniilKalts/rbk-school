package history

import (
	"strings"
	"time"

	"github.com/google/uuid"

	domaincity "github.com/DaniilKalts/rbk-school/7-week/api-service/internal/domain/city"
)

type History struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	City        string
	Temperature float64
	Description string
	RequestedAt time.Time
}

func NewHistory(userID uuid.UUID, city string, temperature float64, description string) (*History, error) {
	h := &History{
		ID:          uuid.New(),
		UserID:      userID,
		City:        domaincity.NormalizeCityName(city),
		Temperature: temperature,
		Description: strings.TrimSpace(description),
	}

	if h.UserID == uuid.Nil {
		return nil, ErrInvalidUserID
	}
	if h.City == "" {
		return nil, ErrInvalidCity
	}

	return h, nil
}
