package queue

import (
	"log"

	"github.com/rapid-downloader/rapid/entry"
	"github.com/rapid-downloader/rapid/setting"
)

type (
	Queue interface {
		Push(entry entry.Entry)
		Pop() entry.Entry
		Len() int
		IsEmpty() bool
	}

	QueueFactory func(setting.Setting) Queue
)

var queueMap = make(map[string]QueueFactory)

func New(provider string, setting setting.Setting) Queue {
	queue, ok := queueMap[provider]
	if !ok {
		log.Panicf("Provider %s is not implemented", provider)
		return nil
	}

	return queue(setting)
}

func RegisterQueue(name string, queue QueueFactory) {
	queueMap[name] = queue
}
