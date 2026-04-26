package history

import (
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
