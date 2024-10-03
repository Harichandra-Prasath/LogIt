package LogIt

import (
	"container/list"
	"fmt"
	"sync"
)

type LogQueue struct {
	lock  sync.RWMutex
	queue *list.List
}

func newLogQueue() *LogQueue {
	return &LogQueue{
		queue: list.New(),
	}
}

func (q *LogQueue) Push(rc Record) {
	q.lock.Lock()
	q.queue.PushBack(rc)
	q.lock.Unlock()
}

func (q *LogQueue) Pop() {
	q.lock.Lock()
	q.queue.Remove(q.queue.Front())
	q.lock.Unlock()
}

func (q *LogQueue) Top() (Record, error) {
	q.lock.RLock()
	r := q.queue.Front()
	if r == nil {
		return Record{}, fmt.Errorf("_empty")
	}
	q.lock.RUnlock()
	return r.Value.(Record), nil
}
