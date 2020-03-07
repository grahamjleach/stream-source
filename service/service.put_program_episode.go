package service

import (
	"time"

	"github.com/grahamjleach/stream-source/domain"
)

func (s *scheduleService) PutProgramEpisode(programID, name string, airDate time.Time, recordingPath string) (episode *domain.ProgramEpisode, err error) {
	episode = &domain.ProgramEpisode{
		ID:            generateID(),
		ProgramID:     programID,
		Name:          name,
		AirDate:       airDate,
		RecordingPath: recordingPath,
	}

	if err = s.schedule.AddProgramEpisode(episode); err != nil {
		return
	}

	if err = s.store.PutProgramEpisode(episode); err != nil {
		//restore
		return
	}

	return
}
