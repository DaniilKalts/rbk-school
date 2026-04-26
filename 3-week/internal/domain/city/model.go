package city

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type City struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	CreatedAt time.Time
}

func New(id uuid.UUID, userID uuid.UUID, name string) (*City, error) {
	c := &City{
		ID:     id,
		UserID: userID,
		Name:   strings.TrimSpace(name),
	}

	if c.ID == uuid.Nil {
		return nil, ErrInvalidID
	}

	if c.UserID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	if c.Name == "" {
		return nil, ErrInvalidName
	}

	return c, nil
}
