package entry

import "container/list"

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
		entries map[string]Entry
	}

	queue struct {
		entries list.List
	}

	Listing struct {
		List
		Queue
	}
)

func NewListing() *Listing {
	return &Listing{
		NewList(),
		NewQueue(),
	}
}

func NewList() List {
	return &entryList{
		map[string]Entry{},
	}
}

func (l *entryList) Init() error {
	// TODO: fetch the lists into persistent disk
	return nil
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

func (l *entryList) Close() error {
	// TODO: save the lists into persistent disk
	return nil
}

func NewQueue() Queue {
	return &queue{}
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
