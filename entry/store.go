package entry

import (
	"github.com/goccy/go-json"
	"github.com/rapid-downloader/rapid/logger"
	"github.com/rapid-downloader/rapid/setting"
	"go.etcd.io/bbolt"
)

type Store interface {
	Get(prefix string, id string) Entry
	GetAll(prefix string) []Entry
	Set(prefix string, id string, val Entry) error
	SetBatch(prefix string, id []string, val []Entry) error
	Delete(prefix string, id string) error
	DeleteAll(prefix string) error
}

type store struct {
	s   *setting.Setting
	log logger.Logger
	db  *bbolt.DB
}

func NewStore(db *bbolt.DB, s *setting.Setting, log logger.Logger) Store {
	return &store{
		s:   s,
		log: log,
		db:  db,
	}
}

func (s *store) Get(prefix string, id string) Entry {
	var out Entry

	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(prefix))
		if err != nil {
			s.log.Print("Error creating bucket on Get:", err.Error())
			return err
		}

		var val entry
		if err := json.Unmarshal(bucket.Get([]byte(id)), &val); err != nil {
			s.log.Print("Error unmarshalling entry:", err.Error())
			return err
		}

		val.Id = id
		out = &val

		return nil
	})

	if err != nil {
		s.log.Print("Error fetching entries from db:", err.Error())
		return nil
	}

	return out
}

func (s *store) GetAll(prefix string) []Entry {
	var entries []Entry

	err := s.db.Batch(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(prefix))
		if err != nil {
			s.log.Print("Error creating bucket on GetAll:", err.Error())
			return err
		}

		return bucket.ForEach(func(k, v []byte) error {
			var entry entry
			if err := json.Unmarshal(v, &entry); err != nil {
				s.log.Print("Error fetching from bucket on get all:", err.Error())
			}

			entry.Id = string(k)
			entries = append(entries, &entry)
			return nil
		})

	})

	if err != nil {
		s.log.Print("Error fetching entries from db:", err.Error())
		return nil
	}

	return entries
}

func (s *store) Set(prefix string, id string, entry Entry) error {
	return nil
}

func (s *store) SetBatch(prefix string, id []string, entries []Entry) error {
	return s.db.Batch(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(prefix))
		if err != nil {
			s.log.Print("Error creating bucket on SetBatch:", err.Error())
			return err
		}

		for i, entry := range entries {
			val, err := json.Marshal(entry)
			if err != nil {
				s.log.Print("Error marshalling for batch operation:", err.Error())
				continue
			}

			if err := bucket.Put([]byte(id[i]), val); err != nil {
				s.log.Print("Error on put set batch:", err.Error())
			}
		}

		return nil
	})
}

func (s *store) Delete(prefix string, id string) error {
	return s.db.View(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(prefix))
		if err != nil {
			s.log.Print("Error creating bucket on Delete:", err.Error())
			return err
		}
		return bucket.Delete([]byte(id))
	})
}

func (s *store) DeleteAll(prefix string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		return tx.DeleteBucket([]byte(prefix))
	})
}
