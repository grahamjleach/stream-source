package main

import (
	"context"
	"sync"
	"time"

	"github.com/grahamjleach/stream-source/domain"
	transport "github.com/grahamjleach/stream-source/interfaces/server/grpc"
	"github.com/grahamjleach/stream-source/usecases"
)

func NewScheduleFeed(ctx context.Context, client transport.ScheduleClient) *ScheduleFeed {
	mutex := sync.Mutex{}

	return &ScheduleFeed{
		ctx:    ctx,
		cond:   sync.NewCond(&mutex),
		client: client,
	}
}

type ScheduleFeed struct {
	ctx          context.Context
	cond         *sync.Cond
	client       transport.ScheduleClient
	scheduleItem domain.ScheduleItem
}

func (s *ScheduleFeed) GetNextScheduleItem() domain.ScheduleItem {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	s.cond.Wait()

	return s.scheduleItem
}

var backoff = time.Second * 10

func (s *ScheduleFeed) LoadItems() {
	nextOffset := usecases.GetCurrentOffset()

	for {
		var (
			item domain.ScheduleItem
			err  error
		)

		if item, err = s.getScheduleItem(nextOffset); err != nil {
			time.Sleep(backoff)
			continue
		}

		if !s.isRunning(item) {
			if item.StartOffset() > nextOffset {
				if item, err = s.getRandomScheduleItem(nextOffset, item.StartOffset()); err != nil {
					time.Sleep(backoff)
					continue
				}
			}

			usecases.SleepUntilOffset(item.StartOffset())
		}

		s.setScheduleItem(item)

		nextOffset = item.EndOffset()
	}
}

func (s *ScheduleFeed) isRunning(scheduleItem domain.ScheduleItem) bool {
	currentOffset := usecases.GetCurrentOffset()

	return scheduleItem.StartOffset() <= currentOffset && scheduleItem.EndOffset() > currentOffset
}

func (s *ScheduleFeed) setScheduleItem(scheduleItem domain.ScheduleItem) {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	s.scheduleItem = scheduleItem

	s.cond.Broadcast()
}

func (s *ScheduleFeed) getScheduleItem(startOffset time.Duration) (scheduleItem domain.ScheduleItem, err error) {
	parameters := &transport.GetEpisodeRequest{
		Offset: startOffset.String(),
	}

	programEpisode, err := s.client.GetEpisode(s.ctx, parameters)
	if err != nil {
		return
	}

	return newScheduleItemFromMessage(programEpisode)
}

func (s *ScheduleFeed) getRandomScheduleItem(startOffset, endOffset time.Duration) (scheduleItem domain.ScheduleItem, err error) {
	parameters := &transport.GetRandomEpisodeRequest{
		MaxDuration: (endOffset - startOffset).String(),
	}

	programEpisode, err := s.client.GetRandomEpisode(s.ctx, parameters)
	if err != nil {
		return
	}

	if scheduleItem, err = newScheduleItemFromMessage(programEpisode); err != nil {
		return
	}

	scheduleItem.Entry.StartOffset = startOffset
	return
}

func newScheduleItemFromMessage(message *transport.ProgramEpisodeDetails) (scheduleItem domain.ScheduleItem, err error) {
	startOffset, err := time.ParseDuration(message.StartOffset)
	if err != nil {
		return
	}

	duration, err := time.ParseDuration(message.Duration)
	if err != nil {
		return
	}

	airDate, err := time.Parse(transport.TimeFormat, message.AirDate)
	if err != nil {
		return
	}

	scheduleItem = domain.ScheduleItem{
		Program: &domain.Program{
			Name:     message.Name,
			HostName: message.HostName,
			Duration: duration,
		},
		Entry: &domain.ProgramScheduleEntry{
			StartOffset: startOffset,
		},
		Episode: &domain.ProgramEpisode{
			Name:          message.EpisodeName,
			AirDate:       airDate,
			RecordingPath: message.RecordingPath,
		},
	}
	return
}
