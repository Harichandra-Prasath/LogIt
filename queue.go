package LogIt

type logQueue struct {

	// buffered channel that acts FIFO Structure
	queue chan record
}

func newLogQueue() *logQueue {
	return &logQueue{
		queue: make(chan record, 10000),
	}
}

func (q *logQueue) push(rc record) {
	q.queue <- rc
}
