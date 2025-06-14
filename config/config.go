package config

import "sync"

type env struct {
	QuestionsFile string
	SnapshotsFile string
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
			QuestionsFile: "questions.json",
			SnapshotsFile: "snapshots.json",
			PageSize:      5,  // Default page size for pagination
			MaxSnapshots:  30, // Maximum number of snapshots to keep
		}
	})
	return envInstance
}
