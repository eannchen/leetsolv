package fileutil

import (
	"encoding/json"
	"os"
)

type FileUtil interface {
	// Load reads data from the specified file into the provided structure.
	// If the file does not exist, it does not return an error and leaves the data in its zero (empty) state.
	Load(data interface{}, filename string) error
	// Save writes data to a temporary file and atomically replaces the target file.
	// This prevents data loss or corruption if the save operation fails or crashes.
	Save(data interface{}, filename string) error
}

type JSONFileUtil struct{}

func NewJSONFileUtil() *JSONFileUtil {
	return &JSONFileUtil{}
}

func (j *JSONFileUtil) Load(data interface{}, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Return nil for non-existent file
		}
		return err
	}
	defer file.Close()

	// Check if file is empty
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	if fileInfo.Size() == 0 {
		return nil // Leave data in its zero state
	}

	if err := json.NewDecoder(file).Decode(data); err != nil {
		return err
	}
	return nil
}

func (j *JSONFileUtil) Save(data interface{}, filename string) error {
	tempFile, err := os.CreateTemp("", "temp_*.json")
	if err != nil {
		return err
	}

	// Ensure cleanup on any error
	cleanup := func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}
	defer cleanup()

	// Write JSON data
	enc := json.NewEncoder(tempFile)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return err
	}

	// Close temp file before rename
	if err := tempFile.Close(); err != nil {
		return err
	}

	// Atomic replace - disable cleanup since we want to keep the file
	cleanup = func() {} // No-op
	if err := os.Rename(tempFile.Name(), filename); err != nil {
		return err
	}
	return nil
}
