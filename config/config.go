package config

import "sync"

type env struct {
	QuestionsFile string
	SnapshotsFile string
	InfoLogFile   string
	ErrorLogFile  string
	PageSize      int
	MaxSnapshots  int
}

var (
	envInstance *env
	once        sync.Once
)

func Env() *env {
	once.Do(func() {
		envInstance = &env{
			QuestionsFile: "questions.json", // Path for the questions file
			SnapshotsFile: "snapshots.json", // Path for the snapshots file
			InfoLogFile:   "info.log",       // Path for the info log file
			ErrorLogFile:  "error.log",      // Path for the error log file
			PageSize:      5,                // Default page size for pagination
			MaxSnapshots:  30,               // Maximum number of snapshots to keep
		}
	})
	return envInstance
}
