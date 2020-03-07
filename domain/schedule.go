package domain

import (
	"errors"
	"sort"
	"sync"
	"time"
)

var (
	ErrorNoScheduleItems   = errors.New("no schedule items")
	ErrorNoProgramEpisodes = errors.New("no episodes")
	ErrorScheduleConflict  = errors.New("schedule conflict")
	ErrorProgramNotFound   = errors.New("program not found")
	ErrorExpiredEpisode    = errors.New("expired episode")
)

var maxEpisodeAge = time.Hour * 24 * 30

type Schedule struct {
	lock     sync.Mutex
	programs map[string]*Program
	episodes map[string]ProgramEpisodes
	items    ScheduleItems //should be sorted
}

func NewSchedule() *Schedule {
	return &Schedule{
		programs: make(map[string]*Program),
		episodes: make(map[string]ProgramEpisodes),
	}
}

func (s *Schedule) AddProgram(program *Program) {
	s.programs[program.ID] = program

	for _, item := range s.items {
		if item.Program.ID == program.ID {
			item.Program = program
		}
	}
}

func (s *Schedule) AddProgramEpisode(episode *ProgramEpisode) error {
	if time.Since(episode.AirDate) > maxEpisodeAge {
		return ErrorExpiredEpisode
	}

	if _, found := s.programs[episode.ProgramID]; !found {
		return ErrorProgramNotFound
	}

	_, found := s.episodes[episode.ProgramID]
	switch {
	case !found:
		s.episodes[episode.ProgramID] = ProgramEpisodes{episode}
	default:
		s.episodes[episode.ProgramID] = append(s.episodes[episode.ProgramID], episode)
		sort.Sort(s.episodes[episode.ProgramID])
	}

	for _, item := range s.items {
		if item.Program.ID == episode.ProgramID {
			item.Episode = s.episodes[episode.ProgramID][0]
		}
	}

	return nil
}

func (s *Schedule) AddProgramScheduleEntry(entry *ProgramScheduleEntry) error {
	program, found := s.programs[entry.ProgramID]
	if !found {
		return ErrorProgramNotFound
	}

	episodes, found := s.episodes[entry.ProgramID]
	if !found || len(episodes) == 0 {
		return ErrorNoProgramEpisodes
	}

	item := ScheduleItem{
		Program: program,
		Entry:   entry,
		Episode: s.episodes[entry.ProgramID][0],
	}

	for _, i := range s.items {
		if i.Overlaps(item) {
			return ErrorScheduleConflict
		}
	}

	s.items = append(s.items, item)
	sort.Sort(s.items)

	return nil
}

func (s *Schedule) GetScheduleItem(offset time.Duration) (item ScheduleItem, err error) {
	if len(s.items) == 0 {
		err = ErrorNoScheduleItems
		return
	}

	for _, i := range s.items {
		if i.StartOffset() <= offset && i.EndOffset() > offset {
			//is running or starts at offset
			return i, nil
		}

		if i.StartOffset() > offset {
			return i, nil
		}
	}

	// how the fuck did we get here?
	// we wrapped, just return the zero entry

	return s.items[0], nil
}

func (s *Schedule) GetRandomScheduleItem(maxDuration time.Duration) (item ScheduleItem, err error) {
	if len(s.items) == 0 {
		err = ErrorNoScheduleItems
		return
	}

	var episode *ProgramEpisode

	for _, p := range s.programs {
		if p.Duration > maxDuration {
			continue
		}

		episodes, found := s.episodes[p.ID]
		if !found || len(episodes) == 0 {
			continue
		}

		episode = episodes.Random()
	}

	if episode == nil {
		err = ErrorNoScheduleItems
		return
	}

	item = ScheduleItem{
		Program: s.programs[episode.ProgramID],
		Episode: episode,
	}
	return
}
