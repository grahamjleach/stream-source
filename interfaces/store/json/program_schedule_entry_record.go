package json

import (
	"encoding/json"
	"time"

	"github.com/grahamjleach/stream-source/domain"
)

type ProgramScheduleEntryRecord struct {
	ID          string `json:"id"`
	ProgramID   string `json:"program_id"`
	StartOffset string `json:"start_offset"`
}

func NewProgramScheduleEntryRecord(entry *domain.ProgramScheduleEntry) *ProgramScheduleEntryRecord {
	return &ProgramScheduleEntryRecord{
		ID:          entry.ID,
		ProgramID:   entry.ProgramID,
		StartOffset: entry.StartOffset.String(),
	}
}

func MarshalProgramScheduleEntry(entry *domain.ProgramScheduleEntry) (js []byte, err error) {
	record := NewProgramScheduleEntryRecord(entry)

	return json.Marshal(record)
}

func UnmarshalProgramScheduleEntry(js []byte) (entry *domain.ProgramScheduleEntry, err error) {
	record := new(ProgramScheduleEntryRecord)

	if err = json.Unmarshal(js, record); err != nil {
		return
	}

	return record.ToProgramScheduleEntry()
}

func (r *ProgramScheduleEntryRecord) ToProgramScheduleEntry() (entry *domain.ProgramScheduleEntry, err error) {
	startOffset, err := time.ParseDuration(r.StartOffset)
	if err != nil {
		return
	}

	entry = &domain.ProgramScheduleEntry{
		ID:          r.ID,
		ProgramID:   r.ProgramID,
		StartOffset: startOffset,
	}
	return
}

type ProgramScheduleEntryRecords []*ProgramScheduleEntryRecord

func (r ProgramScheduleEntryRecords) ToProgramScheduleEntries() (entries []*domain.ProgramScheduleEntry, err error) {
	entries = make([]*domain.ProgramScheduleEntry, len(r))

	for i, record := range r {
		var entry *domain.ProgramScheduleEntry
		if entry, err = record.ToProgramScheduleEntry(); err != nil {
			return
		}
		entries[i] = entry
	}

	return
}
