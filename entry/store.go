package entry

type Store interface {
	Get(id string) Entry
	GetAll() []Entry
	Set(id string, val Entry) error
	SetBatch(id []string, entries []Entry) error
	Delete(id string) error
	DeleteAll() error
}

type memstore struct {
	list List
}

func Memstore() Store {
	return &memstore{
		list: NewList(),
	}
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
