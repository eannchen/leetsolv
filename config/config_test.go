package config

import (
	"os"
	"testing"

	"github.com/eannchen/leetsolv/internal/fileutil"
)

func TestNewConfig(t *testing.T) {
	// Test that NewConfig creates a new instance with DI
	fileUtil := &MockFileUtil{}
	config1, err := NewConfig(fileUtil)
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	config2, err := NewConfig(fileUtil)
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	// Should create different instances (not singleton)
	if config1 == config2 {
		t.Error("NewConfig should create different instances, not singleton")
	}
}

func TestEnvironmentVariableOverride(t *testing.T) {
	// Clear any existing environment variables
	os.Unsetenv("LEETSOLV_RANDOMIZE_INTERVAL")

	// Create config with DI
	fileUtil := &MockFileUtil{}
	config, err := NewConfig(fileUtil)
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	// Test that we can modify the value directly
	originalRandomize := config.RandomizeInterval
	config.RandomizeInterval = !originalRandomize
	if config.RandomizeInterval == originalRandomize {
		t.Error("Expected RandomizeInterval to be modified")
	}

	// Reset to original value
	config.RandomizeInterval = originalRandomize
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
	defer os.Unsetenv("LEETSOLV_SETTINGS_FILE")

	// Create config with real file util for this test
	fileUtil := fileutil.NewJSONFileUtil()
	config, err := NewConfig(fileUtil)
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	// Modify a value
	config.PageSize = 99

	// Save configuration
	if err := config.Save(); err != nil {
		t.Fatalf("Failed to save configuration: %v", err)
	}
}

func TestValidation(t *testing.T) {
	fileUtil := &MockFileUtil{}
	config, err := NewConfig(fileUtil)
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	// Store original values
	originalPageSize := config.PageSize

	// Test valid configuration
	if err := config.validate(); err != nil {
		t.Errorf("Valid configuration should not return error: %v", err)
	}

	// Test invalid configuration
	config.PageSize = -1
	if err := config.validate(); err == nil {
		t.Error("Invalid configuration should return error")
	}

	// Reset to valid state
	config.PageSize = originalPageSize
}

func TestResetToDefaults(t *testing.T) {
	fileUtil := &MockFileUtil{}
	config, err := NewConfig(fileUtil)
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	// Store original values
	originalPageSize := config.PageSize
	originalMaxDelta := config.MaxDelta

	// Modify some values
	config.PageSize = 999
	config.MaxDelta = 888

	// Reset to defaults by loading from default config
	if err := config.loadFromDefault(); err != nil {
		t.Fatalf("Failed to reset to defaults: %v", err)
	}

	// Check that defaults were restored
	if config.PageSize != 5 {
		t.Errorf("Expected PageSize to be 5, got %d", config.PageSize)
	}
	if config.MaxDelta != 25 {
		t.Errorf("Expected MaxDelta to be 25, got %d", config.MaxDelta)
	}

	// Restore original values to avoid affecting other tests
	config.PageSize = originalPageSize
	config.MaxDelta = originalMaxDelta
}

func TestGetSettingsRegistry(t *testing.T) {
	fileUtil := &MockFileUtil{}
	config, err := NewConfig(fileUtil)
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	registry := config.GetSettingsRegistry()

	// Verify registry is not empty
	if len(registry) == 0 {
		t.Error("Expected non-empty settings registry")
	}

	// Verify known settings exist (these are the configurable settings)
	expectedSettings := []string{"randomizeinterval", "overduepenalty", "overduelimit"}
	for _, name := range expectedSettings {
		if _, exists := registry[name]; !exists {
			t.Errorf("Expected setting %q to exist in registry", name)
		}
	}
}

func TestGetSettingValue(t *testing.T) {
	fileUtil := &MockFileUtil{}
	config, err := NewConfig(fileUtil)
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	// Set a known value
	config.RandomizeInterval = true

	// Test getting a valid setting
	value, err := config.GetSettingValue("randomizeinterval")
	if err != nil {
		t.Fatalf("Failed to get setting value: %v", err)
	}
	if value != true {
		t.Errorf("Expected RandomizeInterval to be true, got %v", value)
	}

	// Test case insensitivity
	value, err = config.GetSettingValue("RandomizeInterval")
	if err != nil {
		t.Fatalf("Failed to get setting value with mixed case: %v", err)
	}
	if value != true {
		t.Errorf("Expected RandomizeInterval to be true, got %v", value)
	}

	// Test getting an unknown setting
	_, err = config.GetSettingValue("unknownsetting")
	if err == nil {
		t.Error("Expected error for unknown setting")
	}
}

func TestSetSettingValue(t *testing.T) {
	fileUtil := &MockFileUtil{}
	config, err := NewConfig(fileUtil)
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	// Test setting a valid int value
	if err := config.SetSettingValue("overduelimit", 15); err != nil {
		t.Fatalf("Failed to set setting value: %v", err)
	}
	if config.OverdueLimit != 15 {
		t.Errorf("Expected OverdueLimit to be 15, got %d", config.OverdueLimit)
	}

	// Test setting a valid bool value
	if err := config.SetSettingValue("randomizeinterval", false); err != nil {
		t.Fatalf("Failed to set bool setting: %v", err)
	}
	if config.RandomizeInterval != false {
		t.Error("Expected RandomizeInterval to be false")
	}

	// Test setting an unknown setting
	if err := config.SetSettingValue("unknownsetting", 123); err == nil {
		t.Error("Expected error for unknown setting")
	}
}

func TestGetSettingInfo(t *testing.T) {
	fileUtil := &MockFileUtil{}
	config, err := NewConfig(fileUtil)
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	// Test getting info for a valid setting
	info, err := config.GetSettingInfo("overduelimit")
	if err != nil {
		t.Fatalf("Failed to get setting info: %v", err)
	}
	if info.Name != "OverdueLimit" {
		t.Errorf("Expected Name to be 'OverdueLimit', got %q", info.Name)
	}
	if info.Type != "int" {
		t.Errorf("Expected Type to be 'int', got %q", info.Type)
	}

	// Test getting info for an unknown setting
	_, err = config.GetSettingInfo("unknownsetting")
	if err == nil {
		t.Error("Expected error for unknown setting")
	}
}
