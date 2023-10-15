package api

import (
	"encoding/json"
	"log"

	"go.etcd.io/bbolt"
)

type Store interface {
	Get(id string) *Download
	GetAll() []Download
	Create(id string, val Download) error
	CreateBatch(id []string, entries []Download) error
	Update(id string, val UpdateDownload) error
	Delete(id string) error
	DeleteAll() error
}

type store struct {
	db     *bbolt.DB
	bucket string
}

func NewStore(bucket string, db *bbolt.DB) Store {
	return &store{
		db:     db,
		bucket: bucket,
	}
}

func (s *store) Get(id string) *Download {
	var out Download

	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(s.bucket))
		if err != nil {
			log.Println("Error creating bucket on Get:", err.Error())
			return err
		}

		var val Download
		if err := json.Unmarshal(bucket.Get([]byte(id)), &val); err != nil {
			log.Println("Error unmarshalling entry:", err.Error())
			return err
		}

		val.ID = id
		out = val

		return nil
	})

	if err != nil {
		log.Println("Error fetching entries from db:", err.Error())
		return nil
	}

	return &out
}

func (s *store) GetAll() []Download {
	var entries []Download

	err := s.db.Batch(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(s.bucket))
		if err != nil {
			log.Println("Error creating bucket on GetAll:", err.Error())
			return err
		}

		return bucket.ForEach(func(k, v []byte) error {
			var entry Download
			if err := json.Unmarshal(v, &entry); err != nil {
				log.Println("Error fetching from bucket on get all:", err.Error())
			}

			entry.ID = string(k)
			entries = append(entries, entry)
			return nil
		})

	})

	if err != nil {
		log.Println("Error fetching entries from db:", err.Error())
		return nil
	}

	return entries
}

func (s *store) Create(id string, entry Download) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(s.bucket))
		if err != nil {
			log.Println("Error creating bucket on SetBatch:", err.Error())
			return err
		}

		val, err := json.Marshal(entry)
		if err != nil {
			log.Println("Error marshalling for put operation:", err.Error())
			return err
		}

		return bucket.Put([]byte(id), val)
	})
}

func (s *store) CreateBatch(id []string, entries []Download) error {
	return s.db.Batch(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(s.bucket))
		if err != nil {
			log.Println("Error creating bucket on SetBatch:", err.Error())
			return err
		}

		for i, entry := range entries {
			val, err := json.Marshal(entry)
			if err != nil {
				log.Println("Error marshalling for batch operation:", err.Error())
				continue
			}

			if err := bucket.Put([]byte(id[i]), val); err != nil {
				log.Println("Error on put set batch:", err.Error())
			}
		}

		return nil
	})
}

func (s *store) Update(id string, val UpdateDownload) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(s.bucket))
		if err != nil {
			log.Println("Error creating bucket on SetBatch:", err.Error())
			return err
		}

		var entry Download
		res := bucket.Get([]byte(id))
		if err := json.Unmarshal(res, &entry); err != nil {
			log.Println("Error unmarshalling for update operation:", err.Error())
			return err
		}

		// Update properties only if they are non-nil in val
		if val.URL != nil {
			entry.URL = *val.URL
		}
		if val.Provider != nil {
			entry.Provider = *val.Provider
		}
		if val.Resumable != nil {
			entry.Resumable = *val.Resumable
		}
		if val.Progress != nil {
			entry.Progress = *val.Progress
		}
		if val.Expired != nil {
			entry.Expired = *val.Expired
		}
		if val.ChunkProgress != nil {
			entry.ChunkProgress = val.ChunkProgress
		}
		if val.TimeLeft != nil {
			entry.TimeLeft = *val.TimeLeft
		}
		if val.Speed != nil {
			entry.Speed = *val.Speed
		}
		if val.Status != nil {
			entry.Status = *val.Status
		}

		val, err := json.Marshal(entry)
		if err != nil {
			log.Println("Error marshalling for update operation:", err.Error())
			return err
		}

		return bucket.Put([]byte(id), val)
	})
}

func (s *store) Delete(id string) error {
	return s.db.View(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(s.bucket))
		if err != nil {
			log.Println("Error creating bucket on Delete:", err.Error())
			return err
		}
		return bucket.Delete([]byte(id))
	})
}

func (s *store) DeleteAll() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		return tx.DeleteBucket([]byte(s.bucket))
	})
}
