package domain

import (
	"math/rand"
	"time"
)

type ProgramEpisode struct {
	ID            string
	ProgramID     string
	Name          string
	AirDate       time.Time
	RecordingPath string
}

type ProgramEpisodes []*ProgramEpisode

func (es ProgramEpisodes) Len() int {
	return len(es)
}

func (es ProgramEpisodes) Less(i, j int) bool {
	return es[i].AirDate.Before(es[j].AirDate)
}

func (es ProgramEpisodes) Swap(i, j int) {
	es[i], es[j] = es[j], es[i]
}

func (es ProgramEpisodes) Random() *ProgramEpisode {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	i := int(r.Int31n(int32(len(es))))

	return es[i]
}
