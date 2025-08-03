package config

import (
	"os"
	"testing"
)

// TestConfig provides a way to override configuration for testing
type TestConfig struct {
	QuestionsFile string
	DeltasFile    string
	InfoLogFile   string
	ErrorLogFile  string
	PageSize      int
	MaxDelta      int
	TopKDue       int
	TopKUpcoming  int
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

	// Clean up temp files after test
	t.Cleanup(func() {
		os.Remove(questionsFile.Name())
		os.Remove(deltasFile.Name())
		os.Remove(infoLogFile.Name())
		os.Remove(errorLogFile.Name())
	})

	return &TestConfig{
		QuestionsFile: questionsFile.Name(),
		DeltasFile:    deltasFile.Name(),
		InfoLogFile:   infoLogFile.Name(),
		ErrorLogFile:  errorLogFile.Name(),
		PageSize:      5,
		MaxDelta:      50,
		TopKDue:       10,
		TopKUpcoming:  10,
	}
}

// GetTestEnv returns the test configuration as an env struct
func (tc *TestConfig) GetTestEnv() *env {
	return &env{
		QuestionsFile: tc.QuestionsFile,
		DeltasFile:    tc.DeltasFile,
		InfoLogFile:   tc.InfoLogFile,
		ErrorLogFile:  tc.ErrorLogFile,
		PageSize:      tc.PageSize,
		MaxDelta:      tc.MaxDelta,
		TopKDue:       tc.TopKDue,
		TopKUpcoming:  tc.TopKUpcoming,
	}
}
