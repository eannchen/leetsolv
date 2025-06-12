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
	return fs.load()
}

func (fs *FileStorage) load() ([]core.Question, error) {
	f, err := os.Open(fs.QuestionsFile)
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
	return fs.save(data)
}

func (fs *FileStorage) save(data []core.Question) error {
	// Load the current state and push it to the undo stack
	currentState, err := fs.load()
	if err != nil {
		return err
	}

	if err := fs.pushUndoState(currentState); err != nil {
		return err
	}

	f, err := os.Create(fs.QuestionsFile)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func (fs *FileStorage) pushUndoState(data []core.Question) error {
	// Open the snapshot file
	f, err := os.OpenFile(fs.SnapshotFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	var snapshot [][]core.Question

	// Decode the snapshot file if it is not empty
	fileInfo, err := f.Stat()
	if err != nil {
		return err
	}
	if fileInfo.Size() > 0 {
		err = json.NewDecoder(f).Decode(&snapshot)
		if err != nil {
			return err
		}
	} else {
		// Initialize snapshot as an empty slice if the file is empty
		snapshot = [][]core.Question{}
	}

	// Make a deep copy of the data to avoid modifying the original
	copiedData := make([]core.Question, len(data))
	copy(copiedData, data)
	snapshot = append(snapshot, copiedData)

	// Truncate the file before writing the updated snapshot
	if err := f.Truncate(0); err != nil {
		return err
	}
	if _, err := f.Seek(0, 0); err != nil {
		return err
	}

	// Write the updated snapshot back to the file
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(snapshot)
}

func (fs *FileStorage) Undo() error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	// Open the snapshot file
	f, err := os.OpenFile(fs.SnapshotFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	var snapshot [][]core.Question

	// Decode the snapshot file if it is not empty
	fileInfo, err := f.Stat()
	if err != nil {
		return err
	}
	if fileInfo.Size() > 0 {
		err = json.NewDecoder(f).Decode(&snapshot)
		if err != nil {
			return err
		}
	} else {
		return errors.New("no actions to undo")
	}

	// Get the last state
	lastState := snapshot[len(snapshot)-1]
	// Remove it from the stack
	snapshot = snapshot[:len(snapshot)-1]

	if err = fs.save(lastState); err != nil {
		return err
	}

	// Truncate the file before writing the updated snapshot
	if err := f.Truncate(0); err != nil {
		return err
	}
	if _, err = f.Seek(0, 0); err != nil {
		return err
	}

	// Write the updated snapshot back to the file
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")

	return enc.Encode(snapshot)
}
