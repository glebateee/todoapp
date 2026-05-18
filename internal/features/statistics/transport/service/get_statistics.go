package statistics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/glebateee/todoapp/internal/core/domain"
	core_errors "github.com/glebateee/todoapp/internal/core/errors"
)

func (s *StatisticsService) GetStatistics(
	ctx context.Context,
	userID *int,
	from *time.Time,
	to *time.Time,
) (domain.Statistics, error) {
	if from != nil && to != nil {
		if to.Before(*from) || to.Equal(*from) {
			return domain.Statistics{}, fmt.Errorf("'to' must be aftef 'from': %w", core_errors.ErrInvalidArgument)
		}
	}
	return domain.Statistics{}, nil
}
