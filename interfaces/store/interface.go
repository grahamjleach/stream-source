package store

import "github.com/grahamjleach/stream-source/domain"

type ScheduleStore interface {
	PutProgram(program *domain.Program) error
	PutProgramEpisode(episode *domain.ProgramEpisode) error
	PutProgramScheduleEntry(entry *domain.ProgramScheduleEntry) error
	GetPrograms() ([]*domain.Program, error)
	GetProgramEpisodes() ([]*domain.ProgramEpisode, error)
	GetProgramScheduleEntries() ([]*domain.ProgramScheduleEntry, error)
}
