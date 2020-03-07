package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/grahamjleach/stream-source/domain"
	"github.com/grahamjleach/stream-source/interfaces/store"
)

var timeFormat = time.RFC3339Nano

type jsonStore struct {
	lock     sync.RWMutex
	filePath string
	data     *ScheduleRecord
}

func New(filePath string) (scheduleStore store.ScheduleStore, err error) {
	scheduleRecord := new(ScheduleRecord)

	if file, fileErr := os.Open(filePath); fileErr == nil {
		var js []byte
		if js, err = ioutil.ReadAll(file);err != nil {
			return
		}

		if err = json.Unmarshal(js, scheduleRecord); err != nil {
			return
		}
	}

	scheduleStore = &jsonStore{
		filePath: filePath,
		data:     scheduleRecord,
	}
	return
}

func (s *jsonStore) save() error {
	js, err := json.MarshalIndent(s.data, "", "\t")
	if err != nil {
		return err
	}

	file, err := os.OpenFile(s.filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(js)
	return err
}

func (s *jsonStore) PutProgram(program *domain.Program) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.data.ProgramRecords = append(s.data.ProgramRecords, NewProgramRecord(program))

	return s.save()
}

func (s *jsonStore) PutProgramScheduleEntry(entry *domain.ProgramScheduleEntry) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.data.ProgramScheduleEntryRecords = append(s.data.ProgramScheduleEntryRecords,NewProgramScheduleEntryRecord(entry))

	return s.save()
}

func (s *jsonStore) PutProgramEpisode(episode *domain.ProgramEpisode) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.data.ProgramEpisodeRecords = append(s.data.ProgramEpisodeRecords, NewProgramEpisodeRecord(episode))

	return s.save()
}

func (s *jsonStore) GetPrograms() ([]*domain.Program, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.data.ProgramRecords.ToPrograms()
}

func (s *jsonStore) GetProgramScheduleEntries() ([]*domain.ProgramScheduleEntry, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.data.ProgramScheduleEntryRecords.ToProgramScheduleEntries()
}

func (s *jsonStore) GetProgramEpisodes() ([]*domain.ProgramEpisode, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.data.ProgramEpisodeRecords.ToProgramEpisodes()
}
