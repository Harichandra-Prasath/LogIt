package main

import (
	"os"

	"github.com/Harichandra-Prasath/LogIt/LogIt"
)

func main() {

	opts := LogIt.LoggerOptions{
		Level: LogIt.LEVEL_DEBUG,
		RecordOptions: LogIt.RecordOptions{
			Spacing:   1,
			Colorfull: true,
			Flags:     LogIt.DATE_FLAG | LogIt.TIME_FLAG,
		},
	}

	logger := LogIt.NewLogger(opts, LogIt.NewTextHandler(os.Stdout, os.Stderr))
	defer logger.Flush()

	logger.Info("Hello", "World")

}
