package config

import (
	"os"
	"testing"
)

// TestConfig provides a way to override configuration for testing
type TestConfig struct {
	QuestionsFile     string
	DeltasFile        string
	InfoLogFile       string
	ErrorLogFile      string
	SettingsFile      string
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

// MockEnv creates a test environment with temporary files
func MockEnv(t *testing.T) *TestConfig {
	// Create temporary files for testing
	questionsFile, err := os.CreateTemp("", "test_questions_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp questions file: %v", err)
	}
	questionsFile.Close()

	deltasFile, err := os.CreateTemp("", "test_deltas_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp deltas file: %v", err)
	}
	deltasFile.Close()

	infoLogFile, err := os.CreateTemp("", "test_info_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp info log file: %v", err)
	}
	infoLogFile.Close()

	errorLogFile, err := os.CreateTemp("", "test_error_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp error log file: %v", err)
	}
	errorLogFile.Close()

	settingsFile, err := os.CreateTemp("", "test_settings_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp settings file: %v", err)
	}
	settingsFile.Close()

	// Clean up temp files after test
	t.Cleanup(func() {
		os.Remove(questionsFile.Name())
		os.Remove(deltasFile.Name())
		os.Remove(infoLogFile.Name())
		os.Remove(errorLogFile.Name())
		os.Remove(settingsFile.Name())
	})

	return &TestConfig{
		QuestionsFile:       questionsFile.Name(),
		DeltasFile:          deltasFile.Name(),
		InfoLogFile:         infoLogFile.Name(),
		ErrorLogFile:        errorLogFile.Name(),
		SettingsFile:        settingsFile.Name(),
		PageSize:            3,     // Smaller for testing
		MaxDelta:            20,    // Smaller for testing
		TopKDue:             5,     // Smaller for testing
		TopKUpcoming:        5,     // Smaller for testing
		OverduePenalty:      false, // Disabled for testing
		OverdueLimit:        3,     // Shorter for testing
		RandomizeInterval:   false, // Disabled for testing consistency
		ImportanceWeight:    1.0,   // Neutral for testing
		OverdueWeight:       0.5,   // Standard value
		FamiliarityWeight:   2.0,   // Standard value
		ReviewPenaltyWeight: -1.0,  // Standard value
		EasePenaltyWeight:   -0.5,  // Standard value
	}
}

// GetTestEnv returns the test configuration as an env struct
func (tc *TestConfig) GetTestEnv() *env {
	return &env{
		QuestionsFile:       tc.QuestionsFile,
		DeltasFile:          tc.DeltasFile,
		InfoLogFile:         tc.InfoLogFile,
		ErrorLogFile:        tc.ErrorLogFile,
		SettingsFile:        tc.SettingsFile,
		PageSize:            tc.PageSize,
		MaxDelta:            tc.MaxDelta,
		TopKDue:             tc.TopKDue,
		TopKUpcoming:        tc.TopKUpcoming,
		OverduePenalty:      tc.OverduePenalty,
		OverdueLimit:        tc.OverdueLimit,
		RandomizeInterval:   tc.RandomizeInterval,
		ImportanceWeight:    tc.ImportanceWeight,
		OverdueWeight:       tc.OverdueWeight,
		FamiliarityWeight:   tc.FamiliarityWeight,
		ReviewPenaltyWeight: tc.ReviewPenaltyWeight,
		EasePenaltyWeight:   tc.EasePenaltyWeight,
	}
}

// SetTestEnvironment sets environment variables for testing
func (tc *TestConfig) SetTestEnvironment() {
	os.Setenv("LEETSOLV_QUESTIONS_FILE", tc.QuestionsFile)
	os.Setenv("LEETSOLV_DELTAS_FILE", tc.DeltasFile)
	os.Setenv("LEETSOLV_INFO_LOG_FILE", tc.InfoLogFile)
	os.Setenv("LEETSOLV_ERROR_LOG_FILE", tc.ErrorLogFile)
	os.Setenv("LEETSOLV_SETTINGS_FILE", tc.SettingsFile)
	os.Setenv("LEETSOLV_PAGE_SIZE", "3")
	os.Setenv("LEETSOLV_MAX_DELTA", "20")
	os.Setenv("LEETSOLV_TOP_K_DUE", "5")
	os.Setenv("LEETSOLV_TOP_K_UPCOMING", "5")
}

// ClearTestEnvironment clears test environment variables
func (tc *TestConfig) ClearTestEnvironment() {
	os.Unsetenv("LEETSOLV_QUESTIONS_FILE")
	os.Unsetenv("LEETSOLV_DELTAS_FILE")
	os.Unsetenv("LEETSOLV_INFO_LOG_FILE")
	os.Unsetenv("LEETSOLV_ERROR_LOG_FILE")
	os.Setenv("LEETSOLV_SETTINGS_FILE", "")
	os.Unsetenv("LEETSOLV_PAGE_SIZE")
	os.Unsetenv("LEETSOLV_MAX_DELTA")
	os.Unsetenv("LEETSOLV_TOP_K_DUE")
	os.Unsetenv("LEETSOLV_TOP_K_UPCOMING")
}
