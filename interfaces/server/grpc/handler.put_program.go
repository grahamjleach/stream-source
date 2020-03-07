package grpc

import (
	"context"
	"time"
)

func (h *scheduleHandler) PutProgram(ctx context.Context, request *PutProgramRequest) (message *Program, err error) {
	duration, err := time.ParseDuration(request.GetDuration())
	if err != nil {
		return
	}

	program, err := h.scheduleService.PutProgram(
		request.GetName(),
		request.GetHostName(),
		duration,
	)
	if err != nil {
		return
	}

	message = newProgramMessage(program)
	return
}
