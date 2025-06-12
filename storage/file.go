package storage

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"leetsolv/core"
)

type Storage interface {
	Load() ([]core.Question, error)
	Save([]core.Question) error
	Undo() error
}

type FileStorage struct {
	QuestionsFile string
	SnapshotFile  string
	mu            sync.Mutex
}

func (fs *FileStorage) Load() ([]core.Question, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	return fs.loadQuestions()
}

func (fs *FileStorage) Save(questions []core.Question) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	currentQuestions, err := fs.loadQuestions()
	if err != nil {
		return err
	}

	snapshot, err := fs.loadSnapshot()
	if err != nil {
		return err
	}

	// Deep copy to avoid reference issues
	copiedQuestions := make([]core.Question, len(currentQuestions))
	copy(copiedQuestions, currentQuestions)

	snapshot = append(snapshot, copiedQuestions)

	if err := fs.saveSnapshot(snapshot); err != nil {
		return err
	}
	return fs.saveQuestions(questions)
}

func (fs *FileStorage) Undo() error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	snapshot, err := fs.loadSnapshot()
	if err != nil {
		return err
	}

	if len(snapshot) == 0 {
		return errors.New("no actions to undo")
	}

	// Get the last state and remove it from snapshot
	lastQuestions := snapshot[len(snapshot)-1]
	snapshot = snapshot[:len(snapshot)-1]

	if err := fs.saveQuestions(lastQuestions); err != nil {
		return err
	}
	return fs.saveSnapshot(snapshot)
}

// Private helper methods

func (fs *FileStorage) loadQuestions() ([]core.Question, error) {
	var questions []core.Question
	err := fs.loadJSONFromFile(&questions, fs.QuestionsFile)
	return questions, err
}

func (fs *FileStorage) saveQuestions(questions []core.Question) error {
	return fs.saveJSONToFile(questions, fs.QuestionsFile)
}

func (fs *FileStorage) loadSnapshot() ([][]core.Question, error) {
	var snapshot [][]core.Question
	err := fs.loadJSONFromFile(&snapshot, fs.SnapshotFile)
	return snapshot, err
}

func (fs *FileStorage) saveSnapshot(snapshot [][]core.Question) error {
	return fs.saveJSONToFile(snapshot, fs.SnapshotFile)
}

// loadFile is a generic helper to load JSON data from a file into the provided data structure.
func (fs *FileStorage) loadJSONFromFile(data interface{}, filename string) error {
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

// Generic helper for atomic file writes
func (fs *FileStorage) saveJSONToFile(data interface{}, filename string) error {
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
	return os.Rename(tempFile.Name(), filename)
}
