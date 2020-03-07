package service

import (
	"time"

	"github.com/grahamjleach/stream-source/domain"
)

func (s *scheduleService) PutProgram(name, hostName string, duration time.Duration) (program *domain.Program, err error) {
	program = &domain.Program{
		ID:       generateID(),
		Name:     name,
		HostName: hostName,
		Duration: duration,
	}

	s.schedule.AddProgram(program)

	if err = s.store.PutProgram(program); err != nil {
		//restore
		return
	}

	return
}
