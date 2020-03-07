package grpc

import (
	"time"

	"github.com/grahamjleach/stream-source/domain"
	"github.com/grahamjleach/stream-source/service"
)

var TimeFormat = time.RFC3339Nano

type scheduleHandler struct {
	scheduleService service.ScheduleService
}

func NewScheduleHandler(scheduleService service.ScheduleService) *scheduleHandler {
	return &scheduleHandler{
		scheduleService: scheduleService,
	}
}

func newProgramEpisodeDetailsMessage(scheduleItem domain.ScheduleItem) *ProgramEpisodeDetails {
	program := scheduleItem.GetProgram()
	entry := scheduleItem.GetEntry()
	episode := scheduleItem.GetEpisode()

	return &ProgramEpisodeDetails{
		Name:          program.Name,
		HostName:      program.HostName,
		EpisodeName:   episode.Name,
		StartOffset:   entry.StartOffset.String(),
		Duration:      program.Duration.String(),
		AirDate:       episode.AirDate.Format(TimeFormat),
		RecordingPath: episode.RecordingPath,
	}
}

func newProgramMessage(program *domain.Program) *Program {
	return &Program{
		Id:       program.ID,
		Name:     program.Name,
		HostName: program.HostName,
		Duration: program.Duration.String(),
	}
}

func newProgramScheduleEntryMessage(entry *domain.ProgramScheduleEntry) *ProgramScheduleEntry {
	hours := int(entry.StartOffset.Hours())

	return &ProgramScheduleEntry{
		Id:        entry.ID,
		ProgramId: entry.ProgramID,
		Weekday:   int32(hours / 24),
		Hour:      int32(hours % 24),
	}
}

func newProgramEpisodeMessage(episode *domain.ProgramEpisode) *ProgramEpisode {
	return &ProgramEpisode{
		Id:            episode.ID,
		ProgramId:     episode.ProgramID,
		Name:          episode.Name,
		AirDate:       episode.AirDate.Format(TimeFormat),
		RecordingPath: episode.RecordingPath,
	}
}
