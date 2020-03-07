package domain

import (
	"fmt"
	"time"
)

type ScheduleItem struct {
	Program *Program
	Entry   *ProgramScheduleEntry
	Episode *ProgramEpisode
}

func (i ScheduleItem) GetTitle() string {
	return fmt.Sprintf(
		"%s - %s, %s",
		i.Program.HostName,
		i.Program.Name,
		i.Episode.Name,
	)
}

func (i ScheduleItem) GetEpisodeName() string {
	return i.GetProgram().Name
}

func (i ScheduleItem) GetProgram() *Program {
	if i.Program != nil {
		return i.Program
	}

	return new(Program)
}

func (i ScheduleItem) GetEntry() *ProgramScheduleEntry {
	if i.Entry != nil {
		return i.Entry
	}

	return new(ProgramScheduleEntry)
}

func (i ScheduleItem) GetEpisode() *ProgramEpisode {
	if i.Episode != nil {
		return i.Episode
	}

	return new(ProgramEpisode)
}

func (i ScheduleItem) StartOffset() time.Duration {
	return i.Entry.StartOffset
}

func (i ScheduleItem) EndOffset() time.Duration {
	return i.Entry.StartOffset + i.Program.Duration
}

func (i ScheduleItem) Duration() time.Duration {
	return i.Program.Duration
}

func (i ScheduleItem) Name() string {
	return i.Program.Name
}

func (i ScheduleItem) Overlaps(in ScheduleItem) bool {
	switch {
	case in.StartOffset() == i.StartOffset():
		return true
	case in.StartOffset() < i.StartOffset():
		return in.EndOffset() > i.StartOffset()
	case in.StartOffset() > i.StartOffset():
		return i.EndOffset() > in.StartOffset()
	}

	return false
}
