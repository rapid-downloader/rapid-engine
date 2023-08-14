package entry

import (
	"log"
	"testing"
	"time"

	"github.com/rapid-downloader/rapid/setting"
)

var s = setting.Default()
var url = "https://www.sampledocs.in/DownloadFiles/SampleFile?filename=SampleDocs-Test%20PDF%20File%20With%20Dummy%20Data%20For%20Testing&ext=pdf"

func TestQueuePushAndPop(t *testing.T) {
	q := NewQueue(s)

	entry, _ := Fetch(url)
	q.Push(entry)

	e := q.Pop()

	if e.ID() != entry.ID() {
		t.Errorf("Entry is not the same. Expected to be %s, but got %s", entry.ID(), e.ID())
	}
}

func TestQueueLen(t *testing.T) {
	q := NewQueue(s)

	entry, _ := Fetch(url)
	q.Push(entry)

	if q.Len() != 1 {
		t.Error("Entry is not the same. Expected to be 1, but got", q.Len())
	}

	q.Pop()

	if q.Len() != 0 {
		t.Error("Entry is not the same. Expected to be 1, but got", q.Len())
	}
}

func TestSignal(t *testing.T) {
	q := NewQueue(s)
	qs := q.(QueueSignal)

	n := 3
	go func() {
		for i := 0; i < n; i++ {
			qs.Done()
		}
	}()

	m := 0
	done := make(chan struct{}, 1)
	for range qs.Signal() {
		log.Println("signal received")
		log.Println("processing 2 s")
		time.Sleep(2 * time.Second)

		m++

		if m == n {
			done <- struct{}{}
			break
		}
	}

	select {
	case <-done:
		log.Println("signal is working")
		break
	default:
		t.Error("Signal is not working as intented")
	}

	// to utilize this signal,
	// if the payload is url. fetch all the url
	// push into all entry into queue
	// if queue is empty, then call done
	// otherwise, just wait from the downloader to call done
}

// TODO: test download with queue signal
