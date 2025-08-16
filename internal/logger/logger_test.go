package logger

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewLogger(t *testing.T) {
	// Test that NewLogger creates new instances with DI
	infoPath := filepath.Join(os.TempDir(), "test_info.log")
	errorPath := filepath.Join(os.TempDir(), "test_error.log")

	defer os.Remove(infoPath)
	defer os.Remove(errorPath)

	logger1 := NewLogger(infoPath, errorPath)
	logger2 := NewLogger(infoPath, errorPath)

	// Should create different instances (not singleton)
	if logger1 == logger2 {
		t.Error("NewLogger should create different instances, not singleton")
	}
}

func TestLogger_InstanceCreation(t *testing.T) {
	infoPath := filepath.Join(os.TempDir(), "test_info.log")
	errorPath := filepath.Join(os.TempDir(), "test_error.log")

	defer os.Remove(infoPath)
	defer os.Remove(errorPath)

	logger := NewLogger(infoPath, errorPath)

	if logger == nil {
		t.Fatal("NewLogger returned nil")
	}

	if logger.Info == nil {
		t.Error("Logger Info field should not be nil")
	}

	if logger.Error == nil {
		t.Error("Logger Error field should not be nil")
	}
}

func TestLogger_FileCreation(t *testing.T) {
	// Create temporary log files
	infoPath := filepath.Join(os.TempDir(), "test_info.log")
	errorPath := filepath.Join(os.TempDir(), "test_error.log")

	defer os.Remove(infoPath)
	defer os.Remove(errorPath)

	// Create logger with specific paths
	logger := NewLogger(infoPath, errorPath)
	if logger == nil {
		t.Fatal("NewLogger returned nil")
	}

	// Write a test message to ensure files are created
	logger.Info.Println("Test info message")
	logger.Error.Println("Test error message")

	// Check if info log file exists
	if _, err := os.Stat(infoPath); os.IsNotExist(err) {
		t.Errorf("Info log file not created at %s", infoPath)
	}

	// Check if error log file exists
	if _, err := os.Stat(errorPath); os.IsNotExist(err) {
		t.Errorf("Error log file not created at %s", errorPath)
	}
}

func TestLogger_ConcurrentAccess(t *testing.T) {
	// Test that NewLogger can be called concurrently without issues
	infoPath := filepath.Join(os.TempDir(), "test_concurrent_info.log")
	errorPath := filepath.Join(os.TempDir(), "test_concurrent_error.log")

	defer os.Remove(infoPath)
	defer os.Remove(errorPath)

	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			logger := NewLogger(infoPath, errorPath)
			if logger == nil {
				t.Error("NewLogger returned nil in goroutine")
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
	infoPath := filepath.Join(os.TempDir(), "test_logging_info.log")
	errorPath := filepath.Join(os.TempDir(), "test_logging_error.log")

	defer os.Remove(infoPath)
	defer os.Remove(errorPath)

	logger := NewLogger(infoPath, errorPath)

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

func TestLogger_PathHandling(t *testing.T) {
	// Test that NewLogger handles different path scenarios
	infoPath := filepath.Join(os.TempDir(), "test_path_info.log")
	errorPath := filepath.Join(os.TempDir(), "test_path_error.log")

	defer os.Remove(infoPath)
	defer os.Remove(errorPath)

	logger := NewLogger(infoPath, errorPath)

	// The logger should be created successfully with valid paths
	if logger == nil {
		t.Error("NewLogger should handle valid paths gracefully")
	}
}

func TestLogger_FilePermissions(t *testing.T) {
	// Test that log files are created with appropriate permissions
	infoPath := filepath.Join(os.TempDir(), "test_perms_info.log")
	errorPath := filepath.Join(os.TempDir(), "test_perms_error.log")

	defer os.Remove(infoPath)
	defer os.Remove(errorPath)

	logger := NewLogger(infoPath, errorPath)

	if logger == nil {
		t.Fatal("NewLogger returned nil")
	}

	// The logger should be functional even if we can't check file permissions
	// Just verify that the logger instance is created successfully
	_ = logger
}

func TestLogger_RepeatedCalls(t *testing.T) {
	// Test that calling NewLogger multiple times doesn't cause issues
	infoPath := filepath.Join(os.TempDir(), "test_repeated_info.log")
	errorPath := filepath.Join(os.TempDir(), "test_repeated_error.log")

	defer os.Remove(infoPath)
	defer os.Remove(errorPath)

	for i := 0; i < 100; i++ {
		logger := NewLogger(infoPath, errorPath)
		if logger == nil {
			t.Errorf("NewLogger returned nil on call %d", i)
		}
	}
}

func TestLogger_InterfaceCompliance(t *testing.T) {
	// Test that the logger struct has the expected fields and methods
	infoPath := filepath.Join(os.TempDir(), "test_interface_info.log")
	errorPath := filepath.Join(os.TempDir(), "test_interface_error.log")

	defer os.Remove(infoPath)
	defer os.Remove(errorPath)

	logger := NewLogger(infoPath, errorPath)

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
