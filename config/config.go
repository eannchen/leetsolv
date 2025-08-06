package config

import (
	"os"
	"sync"
)

type env struct {
	QuestionsFile string
	DeltasFile    string
	InfoLogFile   string
	ErrorLogFile  string
	PageSize      int
	MaxDelta      int
	TopKDue       int
	TopKUpcoming  int
}

var (
	envInstance *env
	once        sync.Once
)

func Env() *env {
	once.Do(func() {
		// Get file paths from environment variables, with fallbacks to defaults
		questionsFile := getEnvOrDefault("LEETSOLV_QUESTIONS_FILE", "questions.json")
		deltasFile := getEnvOrDefault("LEETSOLV_DELTAS_FILE", "deltas.json")
		infoLogFile := getEnvOrDefault("LEETSOLV_INFO_LOG_FILE", "info.log")
		errorLogFile := getEnvOrDefault("LEETSOLV_ERROR_LOG_FILE", "error.log")

		envInstance = &env{
			QuestionsFile: questionsFile,
			DeltasFile:    deltasFile,
			InfoLogFile:   infoLogFile,
			ErrorLogFile:  errorLogFile,
			PageSize:      5,  // Default page size for pagination
			MaxDelta:      50, // Maximum number of deltas to keep
			TopKDue:       10, // Top-K due questions to show in summary
			TopKUpcoming:  10, // Top-K upcoming questions to show in summary
		}
	})
	return envInstance
}

// getEnvOrDefault returns the environment variable value if set, otherwise returns the default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
