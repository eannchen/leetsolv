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

func (fs *FileStorage) loadQuestions() ([]core.Question, error) {
	var questions []core.Question

	// Open the questions file (read-only access)
	f, err := os.Open(fs.QuestionsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return questions, nil
		}
		return nil, err
	}
	defer f.Close()

	if err = json.NewDecoder(f).Decode(&questions); err != nil {
		return nil, err
	}
	return questions, nil
}

func (fs *FileStorage) Save(questions []core.Question) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	return fs.saveQuestions(questions)
}

func (fs *FileStorage) saveQuestions(questions []core.Question) error {
	// Load the current questions and push it to the undo stack
	curQuestions, err := fs.loadQuestions()
	if err != nil {
		return err
	}
	if err := fs.pushToSnapshot(curQuestions); err != nil {
		return err
	}

	// Create a temporary file in the same directory as the QuestionsFile
	tempFile, err := os.CreateTemp("", "questions_temp_*.json")
	if err != nil {
		return err
	}
	defer func() {
		// Ensure the temporary file is removed if something goes wrong
		tempFile.Close()
		os.Remove(tempFile.Name())
	}()

	// Write the new questions to the temporary file
	enc := json.NewEncoder(tempFile)
	enc.SetIndent("", "  ")
	if err := enc.Encode(questions); err != nil {
		return err
	}

	// Close the temporary file to flush the data to disk
	if err := tempFile.Close(); err != nil {
		return err
	}

	// Replace the original file with the temporary file
	if err := os.Rename(tempFile.Name(), fs.QuestionsFile); err != nil {
		return err
	}

	return nil
}

func (fs *FileStorage) pushToSnapshot(questions []core.Question) error {
	// Open/Create the questions file (allow read/write)
	f, err := os.OpenFile(fs.SnapshotFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	snapshot := [][]core.Question{}

	// Decode the snapshot file if it is not empty
	fileInfo, err := f.Stat()
	if err != nil {
		return err
	}
	if fileInfo.Size() > 0 {
		if err = json.NewDecoder(f).Decode(&snapshot); err != nil {
			return err
		}
	}

	// Make a deep copy of the data to avoid modifying the original
	copiedData := make([]core.Question, len(questions))
	copy(copiedData, questions)
	snapshot = append(snapshot, copiedData)

	// Create a temporary file for the snapshot
	tempFile, err := os.CreateTemp("", "snapshot_temp_*.json")
	if err != nil {
		return err
	}
	defer func() {
		// Ensure the temporary file is removed if something goes wrong
		tempFile.Close()
		os.Remove(tempFile.Name())
	}()

	// Write the updated snapshot to the temporary file
	enc := json.NewEncoder(tempFile)
	enc.SetIndent("", "  ")
	if err := enc.Encode(snapshot); err != nil {
		return err
	}

	// Close the temporary file to flush the data to disk
	if err := tempFile.Close(); err != nil {
		return err
	}

	// Replace the original snapshot file with the temporary file
	if err := os.Rename(tempFile.Name(), fs.SnapshotFile); err != nil {
		return err
	}

	return nil
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

	if err = fs.saveQuestions(lastState); err != nil {
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
