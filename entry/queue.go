package entry

type (
	Queue interface {
		Push(entry Entry) error
		Pop() Entry
		Len() int
		IsEmpty() bool
	}
)
