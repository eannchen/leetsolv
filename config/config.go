package config

import (
	"os"
	"strconv"
	"sync"
)

type env struct {
	QuestionsFile     string
	DeltasFile        string
	InfoLogFile       string
	ErrorLogFile      string
	PageSize          int
	MaxDelta          int
	TopKDue           int
	TopKUpcoming      int
	OverduePenalty    bool
	OverdueLimit      int
	RandomizeInterval bool
	// Scoring formula weights
	ImportanceWeight    float64
	OverdueWeight       float64
	FamiliarityWeight   float64
	ReviewPenaltyWeight float64
	EasePenaltyWeight   float64
}

var (
	envInstance *env
	once        sync.Once
)

func Env() *env {
	once.Do(func() {
		// Get file paths from environment variables, with fallbacks to defaults
		questionsFile := getEnvOrDefault("LEETSOLV_QUESTIONS_FILE", "questions.test.json")
		deltasFile := getEnvOrDefault("LEETSOLV_DELTAS_FILE", "deltas.test.json")
		infoLogFile := getEnvOrDefault("LEETSOLV_INFO_LOG_FILE", "info.test.log")
		errorLogFile := getEnvOrDefault("LEETSOLV_ERROR_LOG_FILE", "error.test.log")

		envInstance = &env{
			QuestionsFile:     questionsFile,
			DeltasFile:        deltasFile,
			InfoLogFile:       infoLogFile,
			ErrorLogFile:      errorLogFile,
			PageSize:          5,     // Default page size for pagination
			MaxDelta:          50,    // Maximum number of deltas to keep
			TopKDue:           10,    // Top-K due questions to show in summary
			TopKUpcoming:      10,    // Top-K upcoming questions to show in summary
			OverduePenalty:    false, // Enable/disable overdue penalty
			OverdueLimit:      7,     // Days after which overdue questions are at risk of penalty
			RandomizeInterval: true,  // Enable/disable randomized interval
			// Scoring formula weights
			ImportanceWeight:    getEnvFloatOrDefault("LEETSOLV_IMPORTANCE_WEIGHT", 1.5),      // Prioritizes designated importance
			OverdueWeight:       getEnvFloatOrDefault("LEETSOLV_OVERDUE_WEIGHT", 0.5),         // Prioritizes items past their due date
			FamiliarityWeight:   getEnvFloatOrDefault("LEETSOLV_FAMILIARITY_WEIGHT", 3.0),     // Prioritizes historically difficult items
			ReviewPenaltyWeight: getEnvFloatOrDefault("LEETSOLV_REVIEW_PENALTY_WEIGHT", -1.5), // De-prioritizes questions seen many times (prevents leeching)
			EasePenaltyWeight:   getEnvFloatOrDefault("LEETSOLV_EASE_PENALTY_WEIGHT", -1.0),   // De-prioritizes "easier" questions to focus on struggles
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

// getEnvFloatOrDefault returns the environment variable value as float64 if set, otherwise returns the default value
func getEnvFloatOrDefault(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			return f
		}
	}
	return defaultValue
}
