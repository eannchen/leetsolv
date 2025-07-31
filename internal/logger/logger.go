package logger

import (
	"io"
	"log"
	"os"
	"sync"

	"leetsolv/config"
)

type logger struct {
	Info  *log.Logger
	Error *log.Logger
}

var (
	loggerInstance *logger
	once           sync.Once
)

func Logger() *logger {
	once.Do(func() {
		env := config.Env()

		// Open the info log file, creating it if it doesn't exist
		infoFile, err := os.OpenFile(env.InfoLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open info log file: %v", err)
		}

		// Open the error log file, creating it if it doesn't exist
		errorFile, err := os.OpenFile(env.ErrorLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open error log file: %v", err)
		}

		// Combine terminal (os.Stdout / os.Stderr) and file
		errorWriter := io.MultiWriter(os.Stderr, errorFile)

		loggerInstance = &logger{
			Info:  log.New(infoFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
			Error: log.New(errorWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		}
	})
	return loggerInstance
}
