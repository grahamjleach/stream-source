package service

import (
	"time"

	"github.com/grahamjleach/stream-source/domain"
)

type ScheduleService interface {
	PutProgram(name, host string, duration time.Duration) (*domain.Program, error)
	PutProgramScheduleEntry(programID string, weekday time.Weekday, startHour int) (*domain.ProgramScheduleEntry, error)
	PutProgramEpisode(programID, name string, airDate time.Time, recordingPath string) (*domain.ProgramEpisode, error)
	GetScheduleItem(offset time.Duration) (domain.ScheduleItem, error)
	GetRandomScheduleItem(maxDuration time.Duration) (domain.ScheduleItem, error)
}
