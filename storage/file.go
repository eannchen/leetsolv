// Package storage implements the file storage for the leetsolv application.
package storage

import (
	"github.com/eannchen/leetsolv/core"
	"github.com/eannchen/leetsolv/internal/fileutil"
	"github.com/eannchen/leetsolv/internal/search"
)

type Storage interface {
	LoadQuestionStore() (*QuestionStore, error)
	SaveQuestionStore(*QuestionStore) error
	LoadDeltas() ([]core.Delta, error)
	SaveDeltas([]core.Delta) error
	DeleteAllData() error
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
	questionsFileName  string
	deltasFileName     string
	file               fileutil.FileUtil
	questionStoreCache *QuestionStore
	deltasCache        []core.Delta
}

func (fs *FileStorage) LoadQuestionStore() (*QuestionStore, error) {
	// Return from cache if available
	if fs.questionStoreCache != nil {
		return fs.questionStoreCache, nil
	}

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

	// Hydrate the trie nodes
	store.URLTrie.Hydrate()
	store.NoteTrie.Hydrate()

	// Update cache
	fs.questionStoreCache = &store

	return &store, nil
}

func (fs *FileStorage) SaveQuestionStore(store *QuestionStore) error {
	err := fs.file.Save(store, fs.questionsFileName)
	if err != nil {
		return err
	}

	// Update cache after successful save
	fs.questionStoreCache = store

	return nil
}

func (fs *FileStorage) LoadDeltas() ([]core.Delta, error) {
	// Return from cache if available
	if fs.deltasCache != nil {
		return fs.deltasCache, nil
	}

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

// DeleteAllData deletes both questions and deltas files, and invalidates cache
func (fs *FileStorage) DeleteAllData() error {
	// Delete questions file
	if err := fs.file.Delete(fs.questionsFileName); err != nil {
		return err
	}

	// Delete deltas file
	if err := fs.file.Delete(fs.deltasFileName); err != nil {
		return err
	}

	// Invalidate cache
	fs.questionStoreCache = nil
	fs.deltasCache = nil

	return nil
}
