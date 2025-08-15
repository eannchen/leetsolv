package storage

import (
	"sync"

	"leetsolv/core"
	"leetsolv/internal/fileutil"
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

func NewFileStorage(questionsFileName, deltasFileName string, file fileutil.FileUtil) *FileStorage {
	return &FileStorage{
		questionsFileName: questionsFileName,
		deltasFileName:    deltasFileName,
		file:              file,
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
	questionsFileName string
	deltasFileName    string
	file              fileutil.FileUtil
	mu                sync.Mutex

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
	err := fs.file.Load(&store, fs.questionsFileName)
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

	err := fs.file.Save(store, fs.questionsFileName)
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
	err := fs.file.Load(&deltas, fs.deltasFileName)
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

	err := fs.file.Save(deltas, fs.deltasFileName)
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
