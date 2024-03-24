package api

import (
	"encoding/json"
	"fmt"

	"github.com/rapid-downloader/rapid/log"

	"go.etcd.io/bbolt"
)

type Store interface {
	Get(id string) *Download
	GetAll(page, limit int) []Download
	Create(id string, val Download) error
	CreateBatch(id []string, entries []Download) error
	Update(id string, val UpdateDownload) error
	BatchUpdate(ids []string, val []UpdateDownload) error
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
		bucket := tx.Bucket([]byte(s.bucket))

		var val Download
		if err := json.Unmarshal(bucket.Get([]byte(id)), &val); err != nil {
			return fmt.Errorf("error unmarshalling entry:%s", err.Error())
		}

		val.ID = id
		out = val

		return nil
	})

	if err != nil {
		log.Println("error fetching entries from db:", err.Error())
		return nil
	}

	return &out
}

func (s *store) GetAll(page, limit int) []Download {
	var entries []Download

	err := s.db.Batch(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(s.bucket))
		if err != nil {
			return fmt.Errorf("error creating bucket on GetAll:%s", err.Error())
		}

		i := 0

		start := (page - 1) * limit
		end := start + limit

		cursor := bucket.Cursor()

		for k, v := cursor.Last(); k != nil; k, v = cursor.Prev() {
			i++
			if i < start {
				continue
			}

			if i == end {
				break
			}

			var entry Download
			if err := json.Unmarshal(v, &entry); err != nil {
				return fmt.Errorf("error fetching from bucket on get all:%s", err.Error())
			}

			entry.ID = string(k)
			entries = append(entries, entry)
		}

		return nil
	})

	if err != nil {
		log.Println("error fetching entries from db:", err.Error())
		return nil
	}

	return entries
}

func (s *store) Create(id string, entry Download) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(s.bucket))
		if err != nil {
			return fmt.Errorf("error creating bucket on SetBatch:%s", err.Error())
		}

		val, err := json.Marshal(entry)
		if err != nil {
			return fmt.Errorf("error marshalling for put operation:%s", err.Error())
		}

		return bucket.Put([]byte(id), val)
	})
}

func (s *store) CreateBatch(id []string, entries []Download) error {
	return s.db.Batch(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(s.bucket))
		if err != nil {
			return fmt.Errorf("error creating bucket on SetBatch:%s", err.Error())
		}

		for i, entry := range entries {
			val, err := json.Marshal(entry)
			if err != nil {
				return fmt.Errorf("error marshalling for batch operation:%s", err.Error())
			}

			if err := bucket.Put([]byte(id[i]), val); err != nil {
				return fmt.Errorf("error on put set batch:%s", err.Error())
			}
		}

		return nil
	})
}

func (s *store) update(bucket *bbolt.Bucket, id string, toUpdate UpdateDownload) error {

	var entry Download
	res := bucket.Get([]byte(id))
	if err := json.Unmarshal(res, &entry); err != nil {
		return fmt.Errorf("error unmarshalling for update operation:%s", err.Error())
	}

	// Update properties only if they are non-nil in val
	if toUpdate.URL != nil {
		entry.URL = *toUpdate.URL
	}
	if toUpdate.Provider != nil {
		entry.Provider = *toUpdate.Provider
	}
	if toUpdate.Resumable != nil {
		entry.Resumable = *toUpdate.Resumable
	}
	if toUpdate.Progress != nil {
		entry.Progress = *toUpdate.Progress
	}
	if toUpdate.Expired != nil {
		entry.Expired = *toUpdate.Expired
	}
	if toUpdate.DownloadedChunks != nil {
		entry.DownloadedChunks = toUpdate.DownloadedChunks
	}
	if toUpdate.TimeLeft != nil {
		entry.TimeLeft = *toUpdate.TimeLeft
	}
	if toUpdate.Speed != nil {
		entry.Speed = *toUpdate.Speed
	}
	if toUpdate.Status != nil {
		entry.Status = *toUpdate.Status
	}

	val, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("error marshalling for update operation:%s", err.Error())
	}

	return bucket.Put([]byte(id), val)
}

func (s *store) Update(id string, val UpdateDownload) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(s.bucket))
		if err != nil {
			return fmt.Errorf("error creating bucket on SetBatch:%s", err.Error())
		}

		return s.update(bucket, id, val)
	})
}

func (s *store) BatchUpdate(ids []string, val []UpdateDownload) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(s.bucket))
		if err != nil {
			return fmt.Errorf("error creating bucket on SetBatch:%s", err.Error())
		}

		for i, id := range ids {
			if err := s.update(bucket, id, val[i]); err != nil {
				return fmt.Errorf("error updating entries on UpdateAll:%s", err.Error())
			}
		}

		return nil
	})
}

func (s *store) Delete(id string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(s.bucket))
		if err != nil {
			return fmt.Errorf("error creating bucket on Delete:%s", err.Error())
		}

		return bucket.Delete([]byte(id))
	})
}

func (s *store) DeleteAll() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		return tx.DeleteBucket([]byte(s.bucket))
	})
}
