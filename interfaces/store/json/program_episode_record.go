package json

import (
	"encoding/json"
	"time"

	"github.com/grahamjleach/stream-source/domain"
)

type ProgramEpisodeRecord struct {
	ID            string `json:"id"`
	ProgramID     string `json:"program_id"`
	Name          string `json:"name"`
	AirDate       string `json:"air_date"`
	RecordingPath string `json:"recording_path"`
}

func NewProgramEpisodeRecord(episode *domain.ProgramEpisode) *ProgramEpisodeRecord {
	return &ProgramEpisodeRecord{
		ID:            episode.ID,
		ProgramID:     episode.ProgramID,
		Name:          episode.Name,
		AirDate:       episode.AirDate.Format(timeFormat),
		RecordingPath: episode.RecordingPath,
	}
}

func MarshalProgramEpisode(episode *domain.ProgramEpisode) (js []byte, err error) {
	record := NewProgramEpisodeRecord(episode)

	return json.Marshal(record)
}

func UnmarshalProgramEpisode(js []byte) (episode *domain.ProgramEpisode, err error) {
	record := new(ProgramEpisodeRecord)

	if err = json.Unmarshal(js, record); err != nil {
		return
	}

	return record.ToProgramEpisode()
}

func (r *ProgramEpisodeRecord) ToProgramEpisode() (episode *domain.ProgramEpisode, err error) {
	airDate, err := time.Parse(timeFormat, r.AirDate)
	if err != nil {
		return
	}

	episode = &domain.ProgramEpisode{
		ID:            r.ID,
		ProgramID:     r.ProgramID,
		Name:          r.Name,
		AirDate:       airDate,
		RecordingPath: r.RecordingPath,
	}
	return
}

type ProgramEpisodeRecords []*ProgramEpisodeRecord

func (r ProgramEpisodeRecords) ToProgramEpisodes() (episodes []*domain.ProgramEpisode, err error) {
	episodes = make([]*domain.ProgramEpisode, len(r))

	for i, record := range r {
		var episode *domain.ProgramEpisode
		if episode, err = record.ToProgramEpisode(); err != nil {
			return
		}
		episodes[i] = episode
	}

	return
}
