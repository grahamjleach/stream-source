package service

import (
	"time"

	"github.com/grahamjleach/stream-source/domain"
)

func (s *scheduleService) GetRandomScheduleItem(maxDuration time.Duration) (domain.ScheduleItem, error) {
	return s.schedule.GetRandomScheduleItem(maxDuration)
}
