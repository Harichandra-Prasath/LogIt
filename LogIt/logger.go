package LogIt

type LogItLevel int

const LEVEL_DEBUG LogItLevel = 0
const LEVEL_INFO LogItLevel = 1
const LEVEL_WARN LogItLevel = 2
const LEVEL_ERROR LogItLevel = 3

// Options used to configure the logger
type LoggerOptions struct {

	// Least level for the logger. Levels below will be ignored
	Level LogItLevel
}

// Core logger to log the records
type Logger struct {
	Options LoggerOptions
	Handler Handler
}

type Record struct {
	// level of the log message
	Level string

	// actual content of the log
	Message string
}

func NewLogger(opts LoggerOptions, handler Handler) *Logger {

	return &Logger{
		Options: opts,
		Handler: handler,
	}

}

func (l *Logger) Info(message string) {

	// Ignore if the Logger level is higher than Info
	if l.Options.Level > LEVEL_INFO {
		return
	}

	// Create the record
	rc := Record{
		Level:   "INFO",
		Message: message,
	}

	l.Handler.Handle(rc)

}
