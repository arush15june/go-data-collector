package logger

import (
	"fmt"
	"log"
	"os"
)

// NewLogger produces a new logger with required name
func NewLogger(name string) *log.Logger {
	logger := log.New(os.Stdout,
		fmt.Sprintf("%s: ", name),
		log.Ldate|log.Ltime|log.LUTC,
	)

	return logger
}
