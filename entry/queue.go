package entry

import (
	"container/list"
	"sync"

	"github.com/rapid-downloader/rapid/service"
	"github.com/rapid-downloader/rapid/setting"
)

type (
	Queue interface {
		Push(entry Entry)
		Pop() Entry
		Len() int
		IsEmpty() bool
	}
	QueueSignal interface {
		Signal() <-chan struct{}
		Done()
	}

	queue struct {
		list *list.List
		*sync.Mutex
		done chan struct{}
	}
)

func NewQueue(s setting.Setting) Queue {
	return &queue{
		list.New(),
		&sync.Mutex{},
		make(chan struct{}, 1),
	}
}

func (q *queue) Push(entry Entry) {
	q.Lock()
	defer q.Unlock()

	q.list.PushBack(entry)
}

func (q *queue) Pop() Entry {
	q.Lock()
	defer q.Unlock()

	el := q.list.Front()
	if entry, ok := el.Value.(Entry); ok {
		q.list.Remove(el)
		return entry
	}

	return nil
}

func (q *queue) Len() int {
	q.Lock()
	defer q.Unlock()

	return q.list.Len()
}

func (q *queue) IsEmpty() bool {
	q.Lock()
	defer q.Unlock()

	return q.list.Len() == 0
}

func (q *queue) Signal() <-chan struct{} {
	return q.done
}

func (q *queue) Done() {
	q.done <- struct{}{}
}

type (
	queueRunner struct {
		queue Queue
	}

	queueRunnerOption struct {
		downloadProvider string
	}

	QueueRunnerOptions func(o *queueRunnerOption)
)

func UseDownloadProvider(provider string) QueueRunnerOptions {
	return func(o *queueRunnerOption) {
		o.downloadProvider = provider
	}
}

func QueueRunner(queue Queue, options ...QueueRunnerOptions) service.Runner {
	return &queueRunner{
		queue,
	}
}

func (r *queueRunner) Run() error {
	qs := r.queue.(QueueSignal)

	for range qs.Signal() {
		//TODO: download the entry
	}

	return nil
}

func (r *queueRunner) Close() error {
	//TODO: store the queue if any

	return nil
}
