package config

import (
	"os"
	"testing"
)

func TestEnvSingleton(t *testing.T) {
	// Test that Env() returns the same instance
	env1 := Env()
	env2 := Env()
	if env1 != env2 {
		t.Error("Env() should return the same instance")
	}
}

func TestEnvironmentVariableOverride(t *testing.T) {
	// Clear any existing environment variables
	os.Unsetenv("LEETSOLV_PAGE_SIZE")

	// Get the current instance
	env := Env()
	originalPageSize := env.PageSize

	// Test that we can modify the value directly
	env.PageSize = 25
	if env.PageSize != 25 {
		t.Errorf("Expected PageSize to be 25, got %d", env.PageSize)
	}

	// Reset to original value
	env.PageSize = originalPageSize
}

func TestSaveAndLoad(t *testing.T) {
	// Create a temporary settings file
	tempFile, err := os.CreateTemp("", "test_settings_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// Set environment variable to use temp file
	os.Setenv("LEETSOLV_SETTINGS_FILE", tempFile.Name())

	// Get instance and modify a value
	env := Env()
	env.PageSize = 99

	// Save configuration
	if err := env.Save(); err != nil {
		t.Fatalf("Failed to save configuration: %v", err)
	}

	// Clean up
	os.Unsetenv("LEETSOLV_SETTINGS_FILE")
}

func TestValidation(t *testing.T) {
	env := Env()

	// Store original values
	originalPageSize := env.PageSize

	// Test valid configuration
	if err := env.Validate(); err != nil {
		t.Errorf("Valid configuration should not return error: %v", err)
	}

	// Test invalid configuration
	env.PageSize = -1
	if err := env.Validate(); err == nil {
		t.Error("Invalid configuration should return error")
	}

	// Reset to valid state
	env.PageSize = originalPageSize
}

func TestResetToDefaults(t *testing.T) {
	env := Env()

	// Store original values
	originalPageSize := env.PageSize
	originalMaxDelta := env.MaxDelta

	// Modify some values
	env.PageSize = 999
	env.MaxDelta = 888

	// Reset to defaults
	env.ResetToDefaults()

	// Check that defaults were restored
	if env.PageSize != 5 {
		t.Errorf("Expected PageSize to be 5, got %d", env.PageSize)
	}
	if env.MaxDelta != 50 {
		t.Errorf("Expected MaxDelta to be 50, got %d", env.MaxDelta)
	}

	// Restore original values to avoid affecting other tests
	env.PageSize = originalPageSize
	env.MaxDelta = originalMaxDelta
}
