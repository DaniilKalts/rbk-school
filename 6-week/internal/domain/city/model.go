package city

import (
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/google/uuid"
)

type City struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	CreatedAt time.Time
}

func NewCity(userID uuid.UUID, name string) (*City, error) {
	c := &City{
		ID:     uuid.New(),
		UserID: userID,
		Name:   NormalizeCityName(name),
	}

	if c.UserID == uuid.Nil {
		return nil, ErrInvalidUserID
	}
	if c.Name == "" {
		return nil, ErrInvalidName
	}

	return c, nil
}

func NormalizeCityName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return name
	}
	name = strings.ToLower(name)

	r, size := utf8.DecodeRuneInString(name)
	return string(unicode.ToUpper(r)) + name[size:]
}
