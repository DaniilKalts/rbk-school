package weather

import (
	"context"

	"github.com/google/uuid"

	domainhistory "github.com/DaniilKalts/rbk-school/4-week/internal/domain/history"
	domainuser "github.com/DaniilKalts/rbk-school/4-week/internal/domain/user"
	domainweather "github.com/DaniilKalts/rbk-school/4-week/internal/domain/weather"
)

func (s *Service) GetHistory(ctx context.Context, userID uuid.UUID, city string, limit int, offset int) ([]domainhistory.History, error) {
	if userID == uuid.Nil {
		return nil, domainuser.ErrInvalidID
	}

	city = domainweather.NormalizeCityName(city)

	if limit < 0 {
		return nil, domainweather.ErrInvalidLimit
	}
	if offset < 0 {
		return nil, domainweather.ErrInvalidOffset
	}

	_, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.historyRepository.ListHistory(ctx, userID, city, limit, offset)
}
