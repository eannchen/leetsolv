package logger

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInit(t *testing.T) {
	infoPath := filepath.Join(os.TempDir(), "test_init_info.log")
	errorPath := filepath.Join(os.TempDir(), "test_init_error.log")

	defer os.Remove(infoPath)
	defer os.Remove(errorPath)

	// Should not panic
	Init(infoPath, errorPath)

	// Verify loggers are initialized
	if infoLogger == nil {
		t.Error("infoLogger should not be nil after Init")
	}
	if errorLogger == nil {
		t.Error("errorLogger should not be nil after Init")
	}
}

func TestInitNop(t *testing.T) {
	// Should not panic
	InitNop()

	// Verify loggers are initialized (to discard writers)
	if infoLogger == nil {
		t.Error("infoLogger should not be nil after InitNop")
	}
	if errorLogger == nil {
		t.Error("errorLogger should not be nil after InitNop")
	}

	// Should not panic when logging
	Infof("test info message")
	Errorf("test error message")
}

func TestInfof(t *testing.T) {
	infoPath := filepath.Join(os.TempDir(), "test_infof.log")
	errorPath := filepath.Join(os.TempDir(), "test_infof_error.log")

	defer os.Remove(infoPath)
	defer os.Remove(errorPath)

	Init(infoPath, errorPath)

	// Should not panic
	Infof("test message: %s", "hello")

	// Verify the log file was written to
	content, err := os.ReadFile(infoPath)
	if err != nil {
		t.Fatalf("Failed to read info log file: %v", err)
	}

	if len(content) == 0 {
		t.Error("Info log file should not be empty")
	}
}

func TestErrorf(t *testing.T) {
	infoPath := filepath.Join(os.TempDir(), "test_errorf_info.log")
	errorPath := filepath.Join(os.TempDir(), "test_errorf.log")

	defer os.Remove(infoPath)
	defer os.Remove(errorPath)

	Init(infoPath, errorPath)

	// Should not panic
	Errorf("test error: %s", "something went wrong")

	// Verify the log file was written to
	content, err := os.ReadFile(errorPath)
	if err != nil {
		t.Fatalf("Failed to read error log file: %v", err)
	}

	if len(content) == 0 {
		t.Error("Error log file should not be empty")
	}
}

func TestInfof_BeforeInit(t *testing.T) {
	// Reset loggers to nil
	infoLogger = nil
	errorLogger = nil

	// Should not panic when loggers are nil
	Infof("test message")
	Errorf("test error")
}

func TestConcurrentLogging(t *testing.T) {
	infoPath := filepath.Join(os.TempDir(), "test_concurrent_info.log")
	errorPath := filepath.Join(os.TempDir(), "test_concurrent_error.log")

	defer os.Remove(infoPath)
	defer os.Remove(errorPath)

	Init(infoPath, errorPath)

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(n int) {
			Infof("concurrent info message %d", n)
			Errorf("concurrent error message %d", n)
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}
