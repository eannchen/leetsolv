package storage

import (
	"encoding/json"
	"leetsolv/core"
	"os"
	"sync"
)

type Storage interface {
	Load() ([]core.Question, error)
	Save([]core.Question) error
}

type FileStorage struct {
	File string
	mu   sync.Mutex
}

func (fs *FileStorage) Load() ([]core.Question, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	f, err := os.Open(fs.File)
	if err != nil {
		if os.IsNotExist(err) {
			return []core.Question{}, nil
		}
		return nil, err
	}
	defer f.Close()
	var data []core.Question
	err = json.NewDecoder(f).Decode(&data)
	return data, err
}

func (fs *FileStorage) Save(data []core.Question) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	f, err := os.Create(fs.File)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}
