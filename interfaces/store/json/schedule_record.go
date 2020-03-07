package json

type ScheduleRecord struct {
	ProgramRecords              ProgramRecords              `json:"programs,omitempty"`
	ProgramEpisodeRecords       ProgramEpisodeRecords       `json:"program_episodes,omitempty"`
	ProgramScheduleEntryRecords ProgramScheduleEntryRecords `json:"program_schedule_entries,omitempty"`
}
