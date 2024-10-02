package main

import (
	"os"

	"github.com/Harichandra-Prasath/LogIt/LogIt"
)

func main() {

	logger := LogIt.NewLogger(
		LogIt.LoggerOptions{
			Level: 0,
		},
		LogIt.NewDefaultHandler(
			os.Stdout,
		),
	)

	logger.Info("Hello World")
}
