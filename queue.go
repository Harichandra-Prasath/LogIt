package LogIt

import (
	"container/list"
	"fmt"
	"sync"
)

type logQueue struct {
	lock  sync.RWMutex
	queue *list.List
}

func newLogQueue() *logQueue {
	return &logQueue{
		queue: list.New(),
	}
}

func (q *logQueue) push(rc Record) {
	q.lock.Lock()
	q.queue.PushBack(rc)
	q.lock.Unlock()
}

func (q *logQueue) pop() {
	q.lock.Lock()
	q.queue.Remove(q.queue.Front())
	q.lock.Unlock()
}

func (q *logQueue) top() (Record, error) {
	q.lock.RLock()
	r := q.queue.Front()
	if r == nil {
		return Record{}, fmt.Errorf("_empty")
	}
	q.lock.RUnlock()
	return r.Value.(Record), nil
}
