package grpc

import (
	"context"
	"time"
)

func (h *scheduleHandler) PutProgramScheduleEntry(ctx context.Context, request *PutProgramScheduleEntryRequest) (message *ProgramScheduleEntry, err error) {
	entry, err := h.scheduleService.PutProgramScheduleEntry(
		request.GetProgramId(),
		time.Weekday(request.GetWeekday()),
		int(request.GetHour()),
	)
	if err != nil {
		return
	}

	message = newProgramScheduleEntryMessage(entry)
	return
}
