// Package logger implements the logger for the leetsolv application.
package logger

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	Info  *log.Logger
	Error *log.Logger
}

func NewLogger(infoPath, errorPath string) *Logger {
	// Open the info log file, creating it if it doesn't exist
	infoFile, err := os.OpenFile(infoPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open info log file: %v", err)
	}

	// Open the error log file, creating it if it doesn't exist
	errorFile, err := os.OpenFile(errorPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open error log file: %v", err)
	}

	// Combine terminal (os.Stdout / os.Stderr) and file
	errorWriter := io.MultiWriter(os.Stderr, errorFile)

	return &Logger{
		Info:  log.New(infoFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		Error: log.New(errorWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
