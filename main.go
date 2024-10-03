package main

import (
	"github.com/Harichandra-Prasath/LogIt/LogIt"
)

func main() {

	logger := LogIt.DefaultLogger()
	defer logger.Flush()

	logger.Info("Hello", "World")

}
