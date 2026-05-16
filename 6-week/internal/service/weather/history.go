package weather

import (
	"context"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/5-week/internal/domain/history"
	"github.com/DaniilKalts/rbk-school/5-week/internal/domain/user"
	"github.com/DaniilKalts/rbk-school/5-week/internal/domain/weather"

	domaincity "github.com/DaniilKalts/rbk-school/5-week/internal/domain/city"
)

func (s *Service) GetHistory(ctx context.Context, userID uuid.UUID, city string, limit int, offset int) ([]history.History, error) {
	if userID == uuid.Nil {
		return nil, user.ErrInvalidID
	}

	city = domaincity.NormalizeCityName(city)

	if limit < 0 {
		return nil, weather.ErrInvalidLimit
	}
	if offset < 0 {
		return nil, weather.ErrInvalidOffset
	}

	_, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.historyRepository.ListHistory(ctx, userID, city, limit, offset)
}
