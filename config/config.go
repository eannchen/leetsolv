package config

import "sync"

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
		envInstance = &env{
			QuestionsFile: "questions.json", // Path for the questions file
			DeltasFile:    "deltas.json",    // Path for the deltas file
			InfoLogFile:   "info.log",       // Path for the info log file
			ErrorLogFile:  "error.log",      // Path for the error log file
			PageSize:      5,                // Default page size for pagination
			MaxDelta:      50,               // Maximum number of deltas to keep
			TopKDue:       10,               // Top-K due questions to show in summary
			TopKUpcoming:  10,               // Top-K upcoming questions to show in summary
		}
	})
	return envInstance
}
