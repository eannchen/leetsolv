package storage

import (
	"encoding/json"
	"os"
	"sync"

	"leetsolv/core"
	"leetsolv/internal/logger"
	"leetsolv/internal/search"
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
	URLTrie   *search.Trie           `json:"url_trie"`
	NoteTrie  *search.Trie           `json:"note_trie"`
}

type FileStorage struct {
	QuestionsFile string
	DeltasFile    string
	mu            sync.Mutex

	questionStoreCache      *QuestionStore
	questionStoreCacheMutex sync.RWMutex

	deltasCache      []core.Delta
	deltasCacheMutex sync.RWMutex
}

// Lock and Unlock are higher-level locks used to ensure atomicity for a read-and-then-write (lost update or write skew) sequence in the business logic layer.
func (fs *FileStorage) Lock() {
	fs.mu.Lock()
}

func (fs *FileStorage) Unlock() {
	fs.mu.Unlock()
}

func (fs *FileStorage) LoadQuestionStore() (*QuestionStore, error) {
	// Try to load from cache first
	fs.questionStoreCacheMutex.RLock()
	if fs.questionStoreCache != nil {
		fs.questionStoreCacheMutex.RUnlock()
		return fs.questionStoreCache, nil
	}
	fs.questionStoreCacheMutex.RUnlock()

	fs.questionStoreCacheMutex.Lock()
	defer fs.questionStoreCacheMutex.Unlock()

	// Load question store from file
	var store QuestionStore
	err := fs.loadJSONFromFile(&store, fs.QuestionsFile)
	if err != nil {
		return nil, err
	}

	// Initialize empty fields
	if store.Questions == nil {
		store.Questions = make(map[int]*core.Question)
	}
	if store.URLIndex == nil {
		store.URLIndex = make(map[string]int)
	}
	if store.URLTrie == nil {
		store.URLTrie = search.NewTrie(3)
	}
	if store.NoteTrie == nil {
		store.NoteTrie = search.NewTrie(3)
	}

	// Update cache
	fs.questionStoreCache = &store

	return &store, nil
}

func (fs *FileStorage) SaveQuestionStore(store *QuestionStore) error {
	fs.questionStoreCacheMutex.Lock()
	defer fs.questionStoreCacheMutex.Unlock()

	err := fs.saveJSONToFile(store, fs.QuestionsFile)
	if err != nil {
		return err
	}

	// Update cache after successful save
	fs.questionStoreCache = store

	return nil
}

func (fs *FileStorage) LoadDeltas() ([]core.Delta, error) {
	// Try to load from cache first
	fs.deltasCacheMutex.RLock()
	if fs.deltasCache != nil {
		fs.deltasCacheMutex.RUnlock()
		return fs.deltasCache, nil
	}
	fs.deltasCacheMutex.RUnlock()

	fs.deltasCacheMutex.Lock()
	defer fs.deltasCacheMutex.Unlock()

	// Load deltas from file
	var deltas []core.Delta
	err := fs.loadJSONFromFile(&deltas, fs.DeltasFile)
	if err != nil {
		return nil, err
	}

	// Update cache
	fs.deltasCache = deltas

	return deltas, nil
}

func (fs *FileStorage) SaveDeltas(deltas []core.Delta) error {
	fs.deltasCacheMutex.Lock()
	defer fs.deltasCacheMutex.Unlock()

	err := fs.saveJSONToFile(deltas, fs.DeltasFile)
	if err != nil {
		return err
	}

	// Update cache after successful save
	fs.deltasCache = deltas

	return nil
}

// InvalidateCache clears the cache, forcing next load to read from file
func (fs *FileStorage) InvalidateCache() {
	fs.questionStoreCache = nil
	fs.deltasCache = nil
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
