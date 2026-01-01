// Package logger implements the logger for the leetsolv application.
package logger

import (
	"io"
	"log"
	"os"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
)

// Init initializes the package-level loggers with the given file paths.
// Must be called before using Infof or Errorf.
func Init(infoPath, errorPath string) {
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

	// Combine terminal (os.Stderr) and file for errors
	errorWriter := io.MultiWriter(os.Stderr, errorFile)

	infoLogger = log.New(infoFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(errorWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// InitNop initializes loggers that discard all output. Useful for tests.
func InitNop() {
	infoLogger = log.New(io.Discard, "", 0)
	errorLogger = log.New(io.Discard, "", 0)
}

// Infof logs an info message.
func Infof(format string, args ...any) {
	if infoLogger != nil {
		infoLogger.Printf(format, args...)
	}
}

// Errorf logs an error message.
func Errorf(format string, args ...any) {
	if errorLogger != nil {
		errorLogger.Printf(format, args...)
	}
}
