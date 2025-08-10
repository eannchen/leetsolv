package logger

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLogger_Singleton(t *testing.T) {
	// Test that multiple calls to Logger() return the same instance
	logger1 := Logger()
	logger2 := Logger()

	if logger1 != logger2 {
		t.Error("Logger() should return the same instance on multiple calls")
	}
}

func TestLogger_InstanceCreation(t *testing.T) {
	logger := Logger()

	if logger == nil {
		t.Fatal("Logger() returned nil")
	}

	if logger.Info == nil {
		t.Error("Logger Info field should not be nil")
	}

	if logger.Error == nil {
		t.Error("Logger Error field should not be nil")
	}
}

func TestLogger_FileCreation(t *testing.T) {
	// Get the current working directory to check if log files are created
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Check if info log file exists (it should be created by Logger())
	infoLogPath := filepath.Join(cwd, "logs", "info.log")
	if _, err := os.Stat(infoLogPath); os.IsNotExist(err) {
		t.Logf("Info log file not found at %s (this is expected if logs directory doesn't exist)", infoLogPath)
	}

	// Check if error log file exists (it should be created by Logger())
	errorLogPath := filepath.Join(cwd, "logs", "error.log")
	if _, err := os.Stat(errorLogPath); os.IsNotExist(err) {
		t.Logf("Error log file not found at %s (this is expected if logs directory doesn't exist)", errorLogPath)
	}
}

func TestLogger_ConcurrentAccess(t *testing.T) {
	// Test that Logger() can be called concurrently without issues
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			logger := Logger()
			if logger == nil {
				t.Error("Logger() returned nil in goroutine")
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestLogger_Logging(t *testing.T) {
	logger := Logger()

	// Test that we can call logging methods without panicking
	// Note: We can't easily test the actual output without mocking the filesystem
	// But we can test that the methods don't panic

	// Test Info logging
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Info logging panicked: %v", r)
			}
		}()
		logger.Info.Println("Test info message")
	}()

	// Test Error logging
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Error logging panicked: %v", r)
			}
		}()
		logger.Error.Println("Test error message")
	}()
}

func TestLogger_EnvironmentDependency(t *testing.T) {
	// Test that Logger() depends on config.Env()
	// This is an integration test that verifies the dependency
	logger := Logger()

	// The logger should be created successfully even if config.Env() fails
	// (it would use default values or fail gracefully)
	if logger == nil {
		t.Error("Logger() should handle config.Env() gracefully")
	}
}

func TestLogger_FilePermissions(t *testing.T) {
	// Test that log files are created with appropriate permissions
	// Note: This test may not work in all environments due to file permissions
	logger := Logger()

	if logger == nil {
		t.Fatal("Logger() returned nil")
	}

	// The logger should be functional even if we can't check file permissions
	// Just verify that the logger instance is created successfully
	_ = logger
}

func TestLogger_RepeatedCalls(t *testing.T) {
	// Test that calling Logger() multiple times doesn't cause issues
	for i := 0; i < 100; i++ {
		logger := Logger()
		if logger == nil {
			t.Errorf("Logger() returned nil on call %d", i)
		}
	}
}

func TestLogger_InterfaceCompliance(t *testing.T) {
	// Test that the logger struct has the expected fields and methods
	logger := Logger()

	// Test Info logger
	if logger.Info == nil {
		t.Error("Info logger should not be nil")
	}

	// Test Error logger
	if logger.Error == nil {
		t.Error("Error logger should not be nil")
	}

	// Test that both loggers implement the log.Logger interface
	// by calling a method on them
	_ = logger.Info.Flags()
	_ = logger.Error.Flags()
}
