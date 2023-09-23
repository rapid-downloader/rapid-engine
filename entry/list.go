package entry

import (
	"container/list"

	"github.com/rapid-downloader/rapid/db"
	"github.com/rapid-downloader/rapid/logger"
	"github.com/rapid-downloader/rapid/setting"
)

type (
	List interface {
		Insert(entry Entry)
		Remove(id string)
		Find(id string) (Entry, bool)
		Len() int
		IsEmpty() bool
	}

	Queue interface {
		Push(entry Entry)
		Pop() Entry
		Len() int
		IsEmpty() bool
		Range() []Entry
	}

	ListInitter interface {
		Init() error
	}

	ListCloser interface {
		Close() error
	}

	entryList struct {
		s       *setting.Setting
		entries map[string]Entry
		store   Store
		log     logger.Logger
	}

	queue struct {
		s       *setting.Setting
		entries list.List
		store   Store
		log     logger.Logger
	}

	Listing struct {
		List  List
		Queue Queue
	}
)

func NewListing(s *setting.Setting, log logger.Logger) *Listing {
	return &Listing{
		NewList(s, log),
		NewQueue(s, log),
	}
}

func NewList(s *setting.Setting, log logger.Logger) List {
	return &entryList{
		entries: map[string]Entry{},
		store:   NewStore(db.DB(), s, log),
		log:     log,
	}
}

func (l *entryList) populateList() {
	entries := l.store.GetAll("list")

	for _, entry := range entries {
		res, err := Fetch(entry.URL())
		if err != nil {
			l.log.Print("Error fetching init for", entry.Name(), ":", err.Error())
			continue
		}

		l.entries[entry.ID()] = res
	}

	if err := l.store.DeleteAll("list"); err != nil {
		l.log.Print("Error deleting stored data:", err.Error())
	}
}

func (l *entryList) Init() error {
	go l.populateList()
	return nil
}

func (l *entryList) Close() error {
	// TODO: save the lists into persistent disk
	var ids []string
	var entries []Entry

	for k, v := range l.entries {
		ids = append(ids, k)
		entries = append(entries, v)
	}

	return l.store.SetBatch("list", ids, entries)
}

func (l *entryList) Insert(entry Entry) {
	l.entries[entry.ID()] = entry
}

func (l *entryList) Remove(id string) {
	delete(l.entries, id)
}

func (l *entryList) Len() int {
	return len(l.entries)
}

func (l *entryList) Find(id string) (Entry, bool) {
	entry, ok := l.entries[id]
	return entry, ok
}

func (l *entryList) IsEmpty() bool {
	return len(l.entries) == 0
}

func NewQueue(s *setting.Setting, log logger.Logger) Queue {
	return &queue{
		s:     s,
		store: NewStore(db.DB(), s, log),
	}
}

func (q *queue) Init() error {
	// TODO: fetch the queue into persistent disk
	return nil
}

func (q *queue) Push(entry Entry) {
	q.entries.PushBack(entry)
}

func (q *queue) Pop() Entry {
	first := q.entries.Front()
	q.entries.Remove(first)

	return first.Value.(Entry)
}

func (q *queue) Len() int {
	return q.entries.Len()
}

func (q *queue) IsEmpty() bool {
	return q.entries.Len() == 0
}

func (q *queue) Range() []Entry {
	entries := make([]Entry, 0)
	for curr := q.entries.Front(); curr != nil; curr = curr.Next() {
		entries = append(entries, curr.Value.(Entry))
	}

	return entries
}

func (q *queue) Close() error {
	// TODO: save the queue into persistent disk
	return nil
}
