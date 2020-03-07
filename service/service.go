package service

import (
	"github.com/grahamjleach/stream-source/domain"
	"github.com/grahamjleach/stream-source/interfaces/store"
	"github.com/satori/go.uuid"
)

type scheduleService struct {
	store    store.ScheduleStore
	schedule *domain.Schedule
}

func NewScheduleService(scheduleStore store.ScheduleStore) (service *scheduleService, err error) {
	schedule := domain.NewSchedule()

	service = &scheduleService{
		store:    scheduleStore,
		schedule: schedule,
	}

	err = service.hydrateScheduleFromStore()
	return
}

func (s *scheduleService) hydrateScheduleFromStore() error {
	programs, err := s.store.GetPrograms()
	if err != nil {
		return err
	}

	for _, program := range programs {
		s.schedule.AddProgram(program)
	}

	episodes, err := s.store.GetProgramEpisodes()
	if err != nil {
		return err
	}

	for _, episode := range episodes {
		if err := s.schedule.AddProgramEpisode(episode); err != nil {
			switch err {
			case domain.ErrorExpiredEpisode:
				continue
			default:
				return err
			}
		}
	}

	entries, err := s.store.GetProgramScheduleEntries()
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if err := s.schedule.AddProgramScheduleEntry(entry); err != nil {
			return err
		}
	}

	return nil
}

func generateID() string {
	return uuid.NewV4().String()
}
