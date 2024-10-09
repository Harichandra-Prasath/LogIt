package LogIt

import (
	"os"
	"sync"
)

type LogItLevel uint

const LEVEL_DEBUG LogItLevel = 0
const LEVEL_INFO LogItLevel = 1
const LEVEL_WARN LogItLevel = 2
const LEVEL_ERROR LogItLevel = 3

// Options used to configure the logger
type LoggerOptions struct {

	// Least level for the logger. Levels below will be ignored
	Level LogItLevel

	// Options that will be inherited by records for output
	RecordOptions RecordOptions
}

type RecordOptions struct {

	// Flags for the log record
	Flags int

	// Option for Pretty Output
	Colorfull bool

	// Spacing between Flags
	Spacing int
}

// Core logger to log the records
type Logger struct {
	Options  LoggerOptions
	Handler  Handler
	logQueue *logQueue
	stopCh   chan struct{}
	waitgrp  sync.WaitGroup
}

// Core type that represents the final log
type Record struct {
	// level of the log message
	Level string

	//actual content of the message
	Message []string

	Options RecordOptions
}

// Default logger has only StdFlags and Least level as Info;
// Stderr for Error level logs and Stdout for rest
func DefaultLogger() *Logger {
	l := &Logger{
		Options: LoggerOptions{
			Level: LEVEL_INFO,
			RecordOptions: RecordOptions{
				Flags:     STD_FLAG,
				Colorfull: false,
				Spacing:   2,
			},
		},
		Handler: NewTextHandler(
			os.Stdout, os.Stderr,
		),
		logQueue: newLogQueue(),
		stopCh:   make(chan struct{}),
	}
	go l._forward()

	return l
}

// Creates a NewLogger with options and handler
func NewLogger(opts LoggerOptions, handler Handler) *Logger {

	l := &Logger{
		Options:  opts,
		Handler:  handler,
		logQueue: newLogQueue(),
		stopCh:   make(chan struct{}),
	}
	go l._forward()

	return l

}

// Logs the passed message with level INFO
func (l *Logger) Info(message ...string) {

	// Ignore if the Logger level is higher than Info
	if l.Options.Level > LEVEL_INFO {
		return
	}

	l._push("INFO", message...)

}

// Logs the passed message with level DEBUG
func (l *Logger) Debug(message ...string) {

	// Ignore if the Logger level is higher than DEBUG
	if l.Options.Level > LEVEL_DEBUG {
		return
	}

	l._push("DEBUG", message...)

}

// Logs the passed message with level WARN
func (l *Logger) Warn(message ...string) {

	// Ignore if the Logger level is higher than WARN
	if l.Options.Level > LEVEL_WARN {
		return
	}

	l._push("WARN", message...)

}

// Logs the passed message with level ERROR
func (l *Logger) Error(message ...string) {

	// Ignore if the Logger level is higher than ERROR
	if l.Options.Level > LEVEL_ERROR {
		return
	}

	l._push("ERROR", message...)

}

// Push the record to the Queue
func (l *Logger) _push(Level string, message ...string) {

	// Create the record
	rc := Record{
		Level:   Level,
		Message: message,
		Options: l.Options.RecordOptions,
	}

	l.logQueue.push(rc)
}

// Forwards the record to the handler
// FIFO log Processing
func (l *Logger) _forward() {
	l.waitgrp.Add(1)
	defer l.waitgrp.Done()
outer:
	for {
		select {
		case <-l.stopCh:

			// Flush the remaining the records in the Queue
			for len(l.logQueue.queue) > 0 {
				rc := <-l.logQueue.queue
				l.Handler.handle(rc)
			}
			break outer
		case rc := <-l.logQueue.queue:
			l.Handler.handle(rc)
		}

	}
}

// blocking call that flushes all the logs stored in the queue
func (l *Logger) Flush() {
	l.stopCh <- struct{}{}
	l.waitgrp.Wait()
}
