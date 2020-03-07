package service

import (
	"time"

	"github.com/grahamjleach/stream-source/domain"
)

var day = time.Hour * 24

func (s *scheduleService) PutProgramScheduleEntry(programID string, weekday time.Weekday, startHour int) (entry *domain.ProgramScheduleEntry, err error) {
	entry = &domain.ProgramScheduleEntry{
		ID:          generateID(),
		ProgramID:   programID,
		StartOffset: (time.Duration(weekday) * day) + time.Duration(startHour)*time.Hour,
	}

	if err = s.schedule.AddProgramScheduleEntry(entry); err != nil {
		return
	}

	if err = s.store.PutProgramScheduleEntry(entry); err != nil {
		//restore
		return
	}

	return
}
