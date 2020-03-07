package json

import (
	"encoding/json"
	"time"

	"github.com/grahamjleach/stream-source/domain"
)

type ProgramRecord struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	HostName string `json:"host_name"`
	Duration string `json:"duration"`
}

func NewProgramRecord(program *domain.Program) *ProgramRecord {
	return &ProgramRecord{
		ID:       program.ID,
		Name:     program.Name,
		HostName: program.HostName,
		Duration: program.Duration.String(),
	}
}

func MarshalProgram(program *domain.Program) (js []byte, err error) {
	record := NewProgramRecord(program)

	return json.Marshal(record)
}

func UnmarshalProgram(js []byte) (program *domain.Program, err error) {
	record := new(ProgramRecord)

	if err = json.Unmarshal(js, record); err != nil {
		return
	}

	return record.ToProgram()
}

func (r *ProgramRecord) ToProgram() (program *domain.Program, err error) {
	duration, err := time.ParseDuration(r.Duration)
	if err != nil {
		return
	}

	program = &domain.Program{
		ID:       r.ID,
		Name:     r.Name,
		HostName: r.HostName,
		Duration: duration,
	}
	return
}

type ProgramRecords []*ProgramRecord

func (r ProgramRecords) ToPrograms() (programs []*domain.Program, err error) {
	programs = make([]*domain.Program, len(r))

	for i, record := range r {
		var program *domain.Program
		if program, err = record.ToProgram(); err != nil {
			return
		}
		programs[i] = program
	}

	return
}
