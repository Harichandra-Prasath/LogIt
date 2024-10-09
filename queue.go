package LogIt

type logQueue struct {

	// buffered channel that acts FIFO Structure
	queue chan Record
}

func newLogQueue() *logQueue {
	return &logQueue{
		queue: make(chan Record, 10000),
	}
}

func (q *logQueue) push(rc Record) {
	q.queue <- rc
}
