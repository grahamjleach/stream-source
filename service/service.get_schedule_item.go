package service

import (
	"time"

	"github.com/grahamjleach/stream-source/domain"
)

func (s *scheduleService) GetScheduleItem(offset time.Duration) (domain.ScheduleItem, error) {
	return s.schedule.GetScheduleItem(offset)
}
