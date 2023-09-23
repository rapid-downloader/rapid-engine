package entry

import (
	"github.com/goccy/go-json"
	"github.com/rapid-downloader/rapid/db"
	"github.com/rapid-downloader/rapid/log"
	"go.etcd.io/bbolt"
)

type Store interface {
	Get(id string) Entry
	GetAll() []Entry
	Set(id string, val Entry) error
	SetBatch(id []string, entries []Entry) error
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

func (s *store) Get(id string) Entry {
	var out Entry

	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(s.bucket))
		if err != nil {
			log.Println("Error creating bucket on Get:", err.Error())
			return err
		}

		var val entry
		if err := json.Unmarshal(bucket.Get([]byte(id)), &val); err != nil {
			log.Println("Error unmarshalling entry:", err.Error())
			return err
		}

		val.Id = id
		out = &val

		return nil
	})

	if err != nil {
		log.Println("Error fetching entries from db:", err.Error())
		return nil
	}

	return out
}

func (s *store) GetAll() []Entry {
	var entries []Entry

	err := s.db.Batch(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(s.bucket))
		if err != nil {
			log.Println("Error creating bucket on GetAll:", err.Error())
			return err
		}

		return bucket.ForEach(func(k, v []byte) error {
			var entry entry
			if err := json.Unmarshal(v, &entry); err != nil {
				log.Println("Error fetching from bucket on get all:", err.Error())
			}

			entry.Id = string(k)
			entries = append(entries, &entry)
			return nil
		})

	})

	if err != nil {
		log.Println("Error fetching entries from db:", err.Error())
		return nil
	}

	return entries
}

func (s *store) Set(id string, entry Entry) error {
	return nil
}

func (s *store) SetBatch(id []string, entries []Entry) error {
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

type memstore struct {
	disk Store
	list List
}

type MemoryInitter interface {
	Init() error
}
type MemoryCloser interface {
	Close() error
}

func Memstore() Store {
	return &memstore{
		disk: NewStore("list", db.DB()),
		list: NewList(),
	}
}

func (s *memstore) populate() {
	entries := s.disk.GetAll()

	if len(entries) > 0 {
		log.Println("There is unfinished download. Restoring...")
	}

	for _, entry := range entries {
		res, err := Fetch(entry.URL())
		if err != nil {
			log.Println("Error fetching init for", entry.Name(), ":", err.Error())
			continue
		}

		s.list.Insert(entry.ID(), res)
	}

	if err := s.disk.DeleteAll(); err != nil {
		log.Println("Error deleting stored data:", err.Error())
	}
}

func (s *memstore) Init() error {
	go s.populate()
	return nil
}

func (s *memstore) Close() error {
	var ids []string
	var entries []Entry

	for k, v := range s.list.Entries() {
		ids = append(ids, k)
		entries = append(entries, v)
	}

	if len(entries) > 0 {
		log.Println("There is unfinished download. Backuping...")
	}

	return s.disk.SetBatch(ids, entries)
}

func (s *memstore) Get(id string) Entry {
	if entry, ok := s.list.Find(id); ok {
		return entry
	}

	return nil
}

func (s *memstore) GetAll() []Entry {
	entries := make([]Entry, s.list.Len())

	i := 0
	for _, entry := range s.list.Entries() {
		entries[i] = entry
	}

	return entries
}

func (s *memstore) Set(id string, val Entry) error {
	s.list.Insert(id, val)
	return nil
}

func (s *memstore) SetBatch(id []string, entries []Entry) error {
	for i, entry := range entries {
		s.list.Insert(id[i], entry)
	}

	return nil
}

func (s *memstore) Delete(id string) error {
	s.list.Remove(id)
	return nil
}

func (s *memstore) DeleteAll() error {
	entries := s.list.Entries()
	for k := range entries {
		delete(entries, k)
	}

	return nil
}
