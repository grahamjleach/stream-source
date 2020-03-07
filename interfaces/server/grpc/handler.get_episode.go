package grpc

import (
	"context"
	"time"
)

func (h *scheduleHandler) GetEpisode(ctx context.Context, request *GetEpisodeRequest) (message *ProgramEpisodeDetails, err error) {
	offset, err := time.ParseDuration(request.GetOffset())
	if err != nil {
		return
	}

	scheduleItem, err := h.scheduleService.GetScheduleItem(offset)
	if err != nil {
		return
	}

	message = newProgramEpisodeDetailsMessage(scheduleItem)
	return
}
