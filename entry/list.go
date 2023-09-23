package entry

import (
	llist "container/list"
)

type (
	List interface {
		Entries() map[string]Entry
		Insert(id string, entry Entry)
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

	list struct {
		entries map[string]Entry
	}

	queue struct {
		entries *llist.List
	}
)

func NewList() List {
	return &list{
		entries: map[string]Entry{},
	}
}

func (l *list) Insert(id string, entry Entry) {
	l.entries[id] = entry
}

func (l *list) Entries() map[string]Entry {
	return l.entries
}

func (l *list) Remove(id string) {
	delete(l.entries, id)
}

func (l *list) Len() int {
	return len(l.entries)
}

func (l *list) Find(id string) (Entry, bool) {
	entry, ok := l.entries[id]
	return entry, ok
}

func (l *list) IsEmpty() bool {
	return len(l.entries) == 0
}

func NewQueue() Queue {
	return &queue{
		entries: llist.New(),
	}
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
