package grpc

import (
	"context"
	"time"
)

func (h *scheduleHandler) PutProgramEpisode(ctx context.Context, request *PutProgramEpisodeRequest) (message *ProgramEpisode, err error) {
	airDate, err := time.Parse(TimeFormat, request.GetAirDate())
	if err != nil {
		return
	}

	episode, err := h.scheduleService.PutProgramEpisode(
		request.GetProgramId(),
		request.GetName(),
		airDate,
		request.GetRecordingPath(),
	)
	if err != nil {
		return
	}

	message = newProgramEpisodeMessage(episode)
	return
}
