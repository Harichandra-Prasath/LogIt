package main

import (
	"github.com/Harichandra-Prasath/LogIt/LogIt"
)

func main() {

	logger := LogIt.DefaultLogger()
	defer logger.Flush()

	logger.Error("Hello", "World")
	logger.Info("Come")
	logger.Warn("Tellphone")
	logger.Info("dwejoj")
	logger.Warn("JKJ")

}
