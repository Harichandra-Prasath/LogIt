package main

import (
	"os"

	"github.com/Harichandra-Prasath/LogIt/LogIt"
)

func main() {

	logger := LogIt.NewLogger(
		LogIt.LoggerOptions{
			Level: LogIt.LEVEL_DEBUG,
			Flags: LogIt.TIME_FLAG | LogIt.DATE_FLAG,
		},
		LogIt.NewDefaultHandler(
			os.Stdout,
		),
	)

	logger.Info("Hello", "World")
}
