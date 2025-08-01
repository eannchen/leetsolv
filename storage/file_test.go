package storage

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"leetsolv/config"
	"leetsolv/core"
)

// setupTestStorage creates a test storage with temporary files
func setupTestStorage(t *testing.T) (*FileStorage, *config.TestConfig) {
	testConfig := config.MockEnv(t)
	storage := NewFileStorage(testConfig.QuestionsFile, testConfig.DeltasFile)
	return storage, testConfig
}

// createTestQuestion creates a sample question for testing
func createTestQuestion(id int, url string) *core.Question {
	return &core.Question{
		ID:           id,
		URL:          url,
		Note:         "Test question",
		Familiarity:  core.Medium,
		Importance:   core.MediumImportance,
		LastReviewed: time.Now(),
		NextReview:   time.Now().Add(24 * time.Hour),
		ReviewCount:  0,
		EaseFactor:   2.5,
		CreatedAt:    time.Now(),
	}
}

func TestFileStorage_LoadQuestionStore_EmptyFile(t *testing.T) {
	storage, _ := setupTestStorage(t)

	// Test loading from empty file
	store, err := storage.LoadQuestionStore()
	if err != nil {
		t.Fatalf("Expected no error loading empty store, got %v", err)
	}

	if store.MaxID != 0 {
		t.Errorf("Expected MaxID 0, got %d", store.MaxID)
	}

	if len(store.Questions) != 0 {
		t.Errorf("Expected empty questions map, got %d questions", len(store.Questions))
	}

	if len(store.URLIndex) != 0 {
		t.Errorf("Expected empty URL index, got %d entries", len(store.URLIndex))
	}
}

func TestFileStorage_SaveAndLoadQuestionStore(t *testing.T) {
	storage, _ := setupTestStorage(t)

	// Create test data
	question1 := createTestQuestion(1, "https://leetcode.com/problems/test1")
	question2 := createTestQuestion(2, "https://leetcode.com/problems/test2")

	store := &QuestionStore{
		MaxID: 2,
		Questions: map[int]*core.Question{
			1: question1,
			2: question2,
		},
		URLIndex: map[string]int{
			"https://leetcode.com/problems/test1": 1,
			"https://leetcode.com/problems/test2": 2,
		},
	}

	// Save the store
	err := storage.SaveQuestionStore(store)
	if err != nil {
		t.Fatalf("Failed to save question store: %v", err)
	}

	// Load the store
	loadedStore, err := storage.LoadQuestionStore()
	if err != nil {
		t.Fatalf("Failed to load question store: %v", err)
	}

	// Verify the loaded data
	if loadedStore.MaxID != 2 {
		t.Errorf("Expected MaxID 2, got %d", loadedStore.MaxID)
	}

	if len(loadedStore.Questions) != 2 {
		t.Errorf("Expected 2 questions, got %d", len(loadedStore.Questions))
	}

	if len(loadedStore.URLIndex) != 2 {
		t.Errorf("Expected 2 URL index entries, got %d", len(loadedStore.URLIndex))
	}

	// Verify specific questions
	if loadedStore.Questions[1].URL != "https://leetcode.com/problems/test1" {
		t.Errorf("Expected question 1 URL %s, got %s", "https://leetcode.com/problems/test1", loadedStore.Questions[1].URL)
	}

	if loadedStore.Questions[2].URL != "https://leetcode.com/problems/test2" {
		t.Errorf("Expected question 2 URL %s, got %s", "https://leetcode.com/problems/test2", loadedStore.Questions[2].URL)
	}

	// Verify URL index
	if loadedStore.URLIndex["https://leetcode.com/problems/test1"] != 1 {
		t.Errorf("Expected URL index for test1 to be 1, got %d", loadedStore.URLIndex["https://leetcode.com/problems/test1"])
	}

	if loadedStore.URLIndex["https://leetcode.com/problems/test2"] != 2 {
		t.Errorf("Expected URL index for test2 to be 2, got %d", loadedStore.URLIndex["https://leetcode.com/problems/test2"])
	}
}

func TestFileStorage_LoadDeltas_EmptyFile(t *testing.T) {
	storage, _ := setupTestStorage(t)

	// Test loading from empty deltas file
	deltas, err := storage.LoadDeltas()
	if err != nil {
		t.Fatalf("Expected no error loading empty deltas, got %v", err)
	}

	if len(deltas) != 0 {
		t.Errorf("Expected empty deltas slice, got %d deltas", len(deltas))
	}
}

func TestFileStorage_SaveAndLoadDeltas(t *testing.T) {
	storage, _ := setupTestStorage(t)

	// Create test deltas
	question1 := createTestQuestion(1, "https://leetcode.com/problems/test1")
	question2 := createTestQuestion(2, "https://leetcode.com/problems/test2")

	deltas := []core.Delta{
		{
			Action:     core.ActionAdd,
			QuestionID: 1,
			OldState:   nil,
			NewState:   question1,
			CreatedAt:  time.Now(),
		},
		{
			Action:     core.ActionUpdate,
			QuestionID: 1,
			OldState:   question1,
			NewState:   question2,
			CreatedAt:  time.Now(),
		},
	}

	// Save deltas
	err := storage.SaveDeltas(deltas)
	if err != nil {
		t.Fatalf("Failed to save deltas: %v", err)
	}

	// Load deltas
	loadedDeltas, err := storage.LoadDeltas()
	if err != nil {
		t.Fatalf("Failed to load deltas: %v", err)
	}

	// Verify the loaded data
	if len(loadedDeltas) != 2 {
		t.Errorf("Expected 2 deltas, got %d", len(loadedDeltas))
	}

	// Verify first delta
	if loadedDeltas[0].Action != core.ActionAdd {
		t.Errorf("Expected first delta action %s, got %s", core.ActionAdd, loadedDeltas[0].Action)
	}

	if loadedDeltas[0].QuestionID != 1 {
		t.Errorf("Expected first delta question ID 1, got %d", loadedDeltas[0].QuestionID)
	}

	// Verify second delta
	if loadedDeltas[1].Action != core.ActionUpdate {
		t.Errorf("Expected second delta action %s, got %s", core.ActionUpdate, loadedDeltas[1].Action)
	}

	if loadedDeltas[1].QuestionID != 1 {
		t.Errorf("Expected second delta question ID 1, got %d", loadedDeltas[1].QuestionID)
	}
}

func TestFileStorage_ConcurrentAccess(t *testing.T) {
	storage, _ := setupTestStorage(t)

	// Test concurrent access to storage
	done := make(chan bool, 2)

	go func() {
		storage.Lock()
		defer storage.Unlock()
		defer func() { done <- true }()

		// Simulate some work
		time.Sleep(10 * time.Millisecond)
	}()

	go func() {
		storage.Lock()
		defer storage.Unlock()
		defer func() { done <- true }()

		// Simulate some work
		time.Sleep(10 * time.Millisecond)
	}()

	// Wait for both goroutines to complete
	<-done
	<-done

	// If we get here without deadlock, the test passes
}

func TestFileStorage_AtomicWrite(t *testing.T) {
	storage, testConfig := setupTestStorage(t)

	// Create test data
	question := createTestQuestion(1, "https://leetcode.com/problems/test")
	store := &QuestionStore{
		MaxID: 1,
		Questions: map[int]*core.Question{
			1: question,
		},
		URLIndex: map[string]int{
			"https://leetcode.com/problems/test": 1,
		},
	}

	// Save the store
	err := storage.SaveQuestionStore(store)
	if err != nil {
		t.Fatalf("Failed to save question store: %v", err)
	}

	// Verify the file exists and is valid JSON
	file, err := os.Open(testConfig.QuestionsFile)
	if err != nil {
		t.Fatalf("Failed to open saved file: %v", err)
	}
	defer file.Close()

	var loadedStore QuestionStore
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&loadedStore)
	if err != nil {
		t.Fatalf("Failed to decode saved JSON: %v", err)
	}

	// Verify the data was saved correctly
	if loadedStore.MaxID != 1 {
		t.Errorf("Expected MaxID 1, got %d", loadedStore.MaxID)
	}

	if len(loadedStore.Questions) != 1 {
		t.Errorf("Expected 1 question, got %d", len(loadedStore.Questions))
	}
}

func TestFileStorage_LoadQuestionStore_NonExistentFile(t *testing.T) {
	storage, _ := setupTestStorage(t)

	// Remove the file to simulate non-existent file
	os.Remove(storage.QuestionsFile)

	// Should not return an error for non-existent file
	store, err := storage.LoadQuestionStore()
	if err != nil {
		t.Fatalf("Expected no error for non-existent file, got %v", err)
	}

	// Should return empty store
	if store.MaxID != 0 {
		t.Errorf("Expected MaxID 0 for non-existent file, got %d", store.MaxID)
	}

	if len(store.Questions) != 0 {
		t.Errorf("Expected empty questions map for non-existent file, got %d questions", len(store.Questions))
	}
}

func TestFileStorage_LoadDeltas_NonExistentFile(t *testing.T) {
	storage, _ := setupTestStorage(t)

	// Remove the file to simulate non-existent file
	os.Remove(storage.DeltasFile)

	// Should not return an error for non-existent file
	deltas, err := storage.LoadDeltas()
	if err != nil {
		t.Fatalf("Expected no error for non-existent file, got %v", err)
	}

	// Should return empty slice
	if len(deltas) != 0 {
		t.Errorf("Expected empty deltas slice for non-existent file, got %d deltas", len(deltas))
	}
}
