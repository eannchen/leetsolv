package storage

import (
	"encoding/json"
	"os"
	"sync"

	"leetsolv/core"
	"leetsolv/internal/logger"
)

type Storage interface {
	Lock()
	Unlock()
	LoadQuestionStore() (*QuestionStore, error)
	SaveQuestionStore(*QuestionStore) error
	LoadDeltas() ([]core.Delta, error)
	SaveDeltas([]core.Delta) error
}

func NewFileStorage(questionsFile, deltasFile string) *FileStorage {
	return &FileStorage{
		QuestionsFile: questionsFile,
		DeltasFile:    deltasFile,
	}
}

type QuestionStore struct {
	MaxID     int                    `json:"max_id"`
	Questions map[int]*core.Question `json:"questions"`
	URLIndex  map[string]int         `json:"url_index"`
}

type FileStorage struct {
	QuestionsFile string
	DeltasFile    string
	mu            sync.Mutex
}

func (fs *FileStorage) Lock() {
	fs.mu.Lock()
}

func (fs *FileStorage) Unlock() {
	fs.mu.Unlock()
}

func (fs *FileStorage) LoadQuestionStore() (*QuestionStore, error) {
	var jf QuestionStore
	err := fs.loadJSONFromFile(&jf, fs.QuestionsFile)
	if err != nil {
		return nil, err
	}
	if jf.Questions == nil {
		jf.Questions = make(map[int]*core.Question)
	}
	if jf.URLIndex == nil {
		jf.URLIndex = make(map[string]int)
	}
	return &jf, nil
}

func (fs *FileStorage) SaveQuestionStore(store *QuestionStore) error {
	return fs.saveJSONToFile(store, fs.QuestionsFile)
}

func (fs *FileStorage) LoadQuestions() (core.QuestionMap, error) {
	var questions core.QuestionMap
	err := fs.loadJSONFromFile(&questions, fs.QuestionsFile)
	return questions, err
}

func (fs *FileStorage) SaveQuestions(questions core.QuestionMap) error {
	return fs.saveJSONToFile(questions, fs.QuestionsFile)
}

func (fs *FileStorage) LoadDeltas() ([]core.Delta, error) {
	var deltas []core.Delta
	err := fs.loadJSONFromFile(&deltas, fs.DeltasFile)
	return deltas, err
}

func (fs *FileStorage) SaveDeltas(deltas []core.Delta) error {
	return fs.saveJSONToFile(deltas, fs.DeltasFile)
}

// Private helper methods

// loadFile is a generic helper to load JSON data from a file into the provided data structure.
func (fs *FileStorage) loadJSONFromFile(data interface{}, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Return nil for non-existent file
		}
		logger.Logger().Error.Printf("Failed to open file: %s. Error: %v", filename, err)
		return err
	}
	defer file.Close()

	// Check if file is empty
	fileInfo, err := file.Stat()
	if err != nil {
		logger.Logger().Error.Printf("Failed to get file info: %s. Error: %v", filename, err)
		return err
	}
	if fileInfo.Size() == 0 {
		return nil // Leave data in its zero state
	}

	if err := json.NewDecoder(file).Decode(data); err != nil {
		logger.Logger().Error.Printf("Failed to decode JSON from file: %s. Error: %v", filename, err)
		return err
	}
	return nil
}

// Generic helper for atomic file writes
func (fs *FileStorage) saveJSONToFile(data interface{}, filename string) error {
	tempFile, err := os.CreateTemp("", "temp_*.json")
	if err != nil {
		logger.Logger().Error.Printf("Failed to create temporary file for %s. Error: %v", filename, err)
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
		logger.Logger().Error.Printf("Failed to encode JSON to temporary file: %s. Error: %v", tempFile.Name(), err)
		return err
	}

	// Close temp file before rename
	if err := tempFile.Close(); err != nil {
		logger.Logger().Error.Printf("Failed to close temporary file: %s. Error: %v", tempFile.Name(), err)
		return err
	}

	// Atomic replace - disable cleanup since we want to keep the file
	cleanup = func() {} // No-op
	if err := os.Rename(tempFile.Name(), filename); err != nil {
		logger.Logger().Error.Printf("Failed to rename temporary file: %s. Error: %v", tempFile.Name(), err)
		return err
	}
	return nil
}
