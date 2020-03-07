package domain

import "time"

type ProgramScheduleEntry struct {
	ID          string
	ProgramID   string
	StartOffset time.Duration
}
