package bolt

import (
	"github.com/grahamjleach/stream-source/domain"
	"github.com/grahamjleach/stream-source/interfaces/store"
	"github.com/grahamjleach/stream-source/interfaces/store/json"
	"github.com/boltdb/bolt"
)

var (
	programBucketName              = []byte("programs")
	programEpisodeBucketName       = []byte("episodes")
	programScheduleEntryBucketName = []byte("entries")
)

type scheduleStore struct {
	db *bolt.DB
}

func New(databasePath string) (s store.ScheduleStore, err error) {
	db, err := bolt.Open(databasePath, 0600, nil)
	if err != nil {
		return
	}

	s = &scheduleStore{
		db: db,
	}
	return
}

func (s *scheduleStore) PutProgram(program *domain.Program) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(programBucketName)
		if err != nil {
			return err
		}

		js, err := json.MarshalProgram(program)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(program.ID), js)
	})
}

func (s *scheduleStore) PutProgramEpisode(episode *domain.ProgramEpisode) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(programEpisodeBucketName)
		if err != nil {
			return err
		}

		js, err := json.MarshalProgramEpisode(episode)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(episode.ID), js)
	})
}

func (s *scheduleStore) PutProgramScheduleEntry(entry *domain.ProgramScheduleEntry) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(programScheduleEntryBucketName)
		if err != nil {
			return err
		}

		js, err := json.MarshalProgramScheduleEntry(entry)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(entry.ID), js)
	})
}

func (s *scheduleStore) GetPrograms() (programs []*domain.Program, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(programBucketName)
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(_, v []byte) error {
			program, err := json.UnmarshalProgram(v)
			if err != nil {
				return err
			}

			programs = append(programs, program)
			return nil
		})
	})

	return
}

func (s *scheduleStore) GetProgramEpisodes() (episodes []*domain.ProgramEpisode, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(programEpisodeBucketName)
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(_, v []byte) error {
			episode, err := json.UnmarshalProgramEpisode(v)
			if err != nil {
				return err
			}

			episodes = append(episodes, episode)
			return nil
		})
	})

	return
}

func (s *scheduleStore) GetProgramScheduleEntries() (entries []*domain.ProgramScheduleEntry, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(programScheduleEntryBucketName)
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(_, v []byte) error {
			entry, err := json.UnmarshalProgramScheduleEntry(v)
			if err != nil {
				return err
			}

			entries = append(entries, entry)
			return nil
		})
	})

	return
}
