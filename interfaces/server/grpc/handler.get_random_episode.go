package grpc

import (
	"context"
	"time"
)

func (h *scheduleHandler) GetRandomEpisode(ctx context.Context, request *GetRandomEpisodeRequest) (message *ProgramEpisodeDetails, err error) {
	maxDuration, err := time.ParseDuration(request.GetMaxDuration())
	if err != nil {
		return
	}

	scheduleItem, err := h.scheduleService.GetRandomScheduleItem(maxDuration)
	if err != nil {
		return
	}

	message = newProgramEpisodeDetailsMessage(scheduleItem)
	return
}
