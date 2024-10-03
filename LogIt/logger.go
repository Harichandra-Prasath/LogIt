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
	logQueue *LogQueue
	stopCh   chan struct{}
	waitgrp  sync.WaitGroup
}

type Record struct {
	// level of the log message
	Level string

	//actual content of the message
	Message []string

	Options RecordOptions
}

// Default logger has only StdFlags and Least level as Info
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
		Handler: NewLogHandler(
			os.Stdout, os.Stderr,
		),
		logQueue: newLogQueue(),
		stopCh:   make(chan struct{}),
	}
	go l._forward()

	return l
}

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

func (l *Logger) Info(message ...string) {

	// Ignore if the Logger level is higher than Info
	if l.Options.Level > LEVEL_INFO {
		return
	}

	l._push("INFO", message...)

}

func (l *Logger) Debug(message ...string) {

	// Ignore if the Logger level is higher than DEBUG
	if l.Options.Level > LEVEL_DEBUG {
		return
	}

	l._push("DEBUG", message...)

}

func (l *Logger) Warn(message ...string) {

	// Ignore if the Logger level is higher than WARN
	if l.Options.Level > LEVEL_WARN {
		return
	}

	l._push("WARN", message...)

}

func (l *Logger) Error(message ...string) {

	// Ignore if the Logger level is higher than ERROR
	if l.Options.Level > LEVEL_ERROR {
		return
	}

	l._push("ERROR", message...)

}

// Forwards the record to the Hanlder
func (l *Logger) _push(Level string, message ...string) {

	// Create the record
	rc := Record{
		Level:   Level,
		Message: message,
		Options: l.Options.RecordOptions,
	}

	l.logQueue.Push(rc)
}

func (l *Logger) _forward() {
	l.waitgrp.Add(1)
	defer l.waitgrp.Done()
outer:
	for {
		select {
		case <-l.stopCh:
			for {
				rc, err := l.logQueue.Top()
				if err != nil {
					break outer
				}
				l.Handler.handle(rc)
				l.logQueue.Pop()
			}
		default:
			rc, err := l.logQueue.Top()
			if err != nil {
				continue
			}
			l.Handler.handle(rc)
			l.logQueue.Pop()
		}

	}
}

// blocking call that flushes all the logs stored in the queue
func (l *Logger) Flush() {
	l.stopCh <- struct{}{}
	l.waitgrp.Wait()
}
