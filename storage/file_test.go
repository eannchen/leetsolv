package storage

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/eannchen/leetsolv/config"
	"github.com/eannchen/leetsolv/core"
	"github.com/eannchen/leetsolv/internal/fileutil"
	"github.com/eannchen/leetsolv/internal/search"
)

// Fixed test time for deterministic tests
var testTime = time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)

// setupTestStorage creates a test storage with temporary files
func setupTestStorage(t *testing.T) (*FileStorage, *config.TestConfig) {
	testConfig, _ := config.MockEnv(t)
	fileUtil := fileutil.NewJSONFileUtil()
	storage := NewFileStorage(testConfig.QuestionsFile, testConfig.DeltasFile, fileUtil)
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
		LastReviewed: testTime,
		NextReview:   testTime.Add(24 * time.Hour),
		ReviewCount:  0,
		EaseFactor:   2.5,
		CreatedAt:    testTime,
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

	// Verify trie initialization
	if store.URLTrie == nil {
		t.Error("Expected URLTrie to be initialized")
	}
	if store.NoteTrie == nil {
		t.Error("Expected NoteTrie to be initialized")
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
		URLTrie:  search.NewTrie(3),
		NoteTrie: search.NewTrie(3),
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
			CreatedAt:  testTime,
		},
		{
			Action:     core.ActionUpdate,
			QuestionID: 1,
			OldState:   question1,
			NewState:   question2,
			CreatedAt:  testTime,
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
	storage, testConfig := setupTestStorage(t)

	// Remove the file to simulate non-existent file
	os.Remove(testConfig.QuestionsFile)

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
	storage, testConfig := setupTestStorage(t)

	// Remove the file to simulate non-existent file
	os.Remove(testConfig.DeltasFile)

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

// NEW TESTS FOR BETTER BUG DETECTION

func TestFileStorage_CorruptedJSONFile(t *testing.T) {
	storage, testConfig := setupTestStorage(t)

	// Write corrupted JSON to file
	corruptedJSON := `{"max_id": 1, "questions": {"1": {"id": 1, "url": "test"}}` // Missing closing brace
	err := os.WriteFile(testConfig.QuestionsFile, []byte(corruptedJSON), 0644)
	if err != nil {
		t.Fatalf("Failed to write corrupted JSON: %v", err)
	}

	// Should handle corrupted JSON gracefully
	_, err = storage.LoadQuestionStore()
	if err == nil {
		t.Error("Expected error when loading corrupted JSON")
	}
}

func TestFileStorage_EmptyJSONFile(t *testing.T) {
	storage, testConfig := setupTestStorage(t)

	// Write empty file
	err := os.WriteFile(testConfig.QuestionsFile, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to write empty file: %v", err)
	}

	// Should handle empty file gracefully
	store, err := storage.LoadQuestionStore()
	if err != nil {
		t.Fatalf("Expected no error for empty file, got %v", err)
	}

	if store.MaxID != 0 {
		t.Errorf("Expected MaxID 0 for empty file, got %d", store.MaxID)
	}
}

func TestFileStorage_FilePermissionIssues(t *testing.T) {
	_, testConfig := setupTestStorage(t)

	// Test with a directory that doesn't exist (more reliable than read-only permissions)
	nonExistentDir := "/non/existent/directory"
	storageWithBadPath := NewFileStorage(nonExistentDir+"/questions.json", testConfig.DeltasFile, fileutil.NewJSONFileUtil())

	store := &QuestionStore{MaxID: 2}
	err := storageWithBadPath.SaveQuestionStore(store)
	if err == nil {
		t.Error("Expected error when saving to non-existent directory")
	}
}

func TestFileStorage_LargeDataSet(t *testing.T) {
	storage, _ := setupTestStorage(t)

	// Create a large dataset
	store := &QuestionStore{
		MaxID:     1000,
		Questions: make(map[int]*core.Question),
		URLIndex:  make(map[string]int),
		URLTrie:   search.NewTrie(3),
		NoteTrie:  search.NewTrie(3),
	}

	// Add 1000 questions
	for i := 1; i <= 1000; i++ {
		question := createTestQuestion(i, "https://leetcode.com/problems/test")
		store.Questions[i] = question
		store.URLIndex[question.URL] = i
	}

	// Save large dataset
	err := storage.SaveQuestionStore(store)
	if err != nil {
		t.Fatalf("Failed to save large dataset: %v", err)
	}

	// Load large dataset
	loadedStore, err := storage.LoadQuestionStore()
	if err != nil {
		t.Fatalf("Failed to load large dataset: %v", err)
	}

	if len(loadedStore.Questions) != 1000 {
		t.Errorf("Expected 1000 questions, got %d", len(loadedStore.Questions))
	}
}

func TestFileStorage_TrieInitialization(t *testing.T) {
	storage, _ := setupTestStorage(t)

	// Test loading with nil tries
	store := &QuestionStore{
		MaxID:     1,
		Questions: map[int]*core.Question{1: createTestQuestion(1, "test")},
		URLIndex:  map[string]int{"test": 1},
		// URLTrie and NoteTrie are nil
	}

	err := storage.SaveQuestionStore(store)
	if err != nil {
		t.Fatalf("Failed to save store with nil tries: %v", err)
	}

	// Invalidate cache to force reload from file
	storage.InvalidateCache()

	// Load should initialize nil tries
	loadedStore, err := storage.LoadQuestionStore()
	if err != nil {
		t.Fatalf("Failed to load store: %v", err)
	}

	if loadedStore.URLTrie == nil {
		t.Error("Expected URLTrie to be initialized")
	}
	if loadedStore.NoteTrie == nil {
		t.Error("Expected NoteTrie to be initialized")
	}
}

func TestFileStorage_DeltaLimitHandling(t *testing.T) {
	storage, _ := setupTestStorage(t)

	// Create more deltas than the limit
	deltas := make([]core.Delta, 1000)
	for i := 0; i < 1000; i++ {
		deltas[i] = core.Delta{
			Action:     core.ActionAdd,
			QuestionID: i,
			CreatedAt:  testTime,
		}
	}

	err := storage.SaveDeltas(deltas)
	if err != nil {
		t.Fatalf("Failed to save many deltas: %v", err)
	}

	loadedDeltas, err := storage.LoadDeltas()
	if err != nil {
		t.Fatalf("Failed to load deltas: %v", err)
	}

	// Should handle large number of deltas
	if len(loadedDeltas) != 1000 {
		t.Errorf("Expected 1000 deltas, got %d", len(loadedDeltas))
	}
}

func TestFileStorage_InvalidJSONStructure(t *testing.T) {
	storage, testConfig := setupTestStorage(t)

	// Write JSON with wrong structure
	invalidJSON := `{"invalid_field": "value"}`
	err := os.WriteFile(testConfig.QuestionsFile, []byte(invalidJSON), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid JSON: %v", err)
	}

	// Should handle invalid structure gracefully
	store, err := storage.LoadQuestionStore()
	if err != nil {
		t.Fatalf("Expected no error for invalid structure, got %v", err)
	}

	// Should return empty store
	if store.MaxID != 0 {
		t.Errorf("Expected MaxID 0 for invalid structure, got %d", store.MaxID)
	}
}

func TestFileStorage_TempFileCleanup(t *testing.T) {
	storage, _ := setupTestStorage(t)

	// Create a store
	store := &QuestionStore{
		MaxID:     1,
		Questions: map[int]*core.Question{1: createTestQuestion(1, "test")},
	}

	// Save multiple times to test temp file cleanup
	for i := 0; i < 10; i++ {
		err := storage.SaveQuestionStore(store)
		if err != nil {
			t.Fatalf("Failed to save store iteration %d: %v", i, err)
		}
	}

	// Verify the final file is correct
	loadedStore, err := storage.LoadQuestionStore()
	if err != nil {
		t.Fatalf("Failed to load store: %v", err)
	}

	if loadedStore.MaxID != 1 {
		t.Errorf("Expected MaxID 1, got %d", loadedStore.MaxID)
	}
}

// NEW TESTS FOR CACHING FUNCTIONALITY

func TestFileStorage_CacheBehavior(t *testing.T) {
	storage, _ := setupTestStorage(t)

	// First load should populate cache
	store1, err := storage.LoadQuestionStore()
	if err != nil {
		t.Fatalf("Failed to load question store: %v", err)
	}

	// Second load should return cached data (same pointer)
	store2, err := storage.LoadQuestionStore()
	if err != nil {
		t.Fatalf("Failed to load question store second time: %v", err)
	}

	// Should return the same cached instance
	if store1 != store2 {
		t.Error("Expected second load to return cached instance")
	}
}

func TestFileStorage_CacheUpdateOnSave(t *testing.T) {
	storage, _ := setupTestStorage(t)

	// Create and save initial store
	initialStore := &QuestionStore{
		MaxID:     1,
		Questions: map[int]*core.Question{1: createTestQuestion(1, "test1")},
		URLIndex:  map[string]int{"test1": 1},
		URLTrie:   search.NewTrie(3),
		NoteTrie:  search.NewTrie(3),
	}

	err := storage.SaveQuestionStore(initialStore)
	if err != nil {
		t.Fatalf("Failed to save initial store: %v", err)
	}

	// Load should return cached version
	_, err = storage.LoadQuestionStore()
	if err != nil {
		t.Fatalf("Failed to load cached store: %v", err)
	}

	// Modify and save updated store
	updatedStore := &QuestionStore{
		MaxID: 2,
		Questions: map[int]*core.Question{
			1: createTestQuestion(1, "test1"),
			2: createTestQuestion(2, "test2"),
		},
		URLIndex: map[string]int{"test1": 1, "test2": 2},
		URLTrie:  search.NewTrie(3),
		NoteTrie: search.NewTrie(3),
	}

	err = storage.SaveQuestionStore(updatedStore)
	if err != nil {
		t.Fatalf("Failed to save updated store: %v", err)
	}

	// Load should return updated cached version
	loadedStore, err := storage.LoadQuestionStore()
	if err != nil {
		t.Fatalf("Failed to load updated store: %v", err)
	}

	if loadedStore.MaxID != 2 {
		t.Errorf("Expected MaxID 2 after update, got %d", loadedStore.MaxID)
	}

	if len(loadedStore.Questions) != 2 {
		t.Errorf("Expected 2 questions after update, got %d", len(loadedStore.Questions))
	}

	// Should return the same instance as the updated store
	if loadedStore != updatedStore {
		t.Error("Expected load to return the same instance as the updated store")
	}
}

func TestFileStorage_DeltasCacheBehavior(t *testing.T) {
	storage, _ := setupTestStorage(t)

	// First load should populate cache
	deltas1, err := storage.LoadDeltas()
	if err != nil {
		t.Fatalf("Failed to load deltas: %v", err)
	}

	// Second load should return cached data (same slice)
	deltas2, err := storage.LoadDeltas()
	if err != nil {
		t.Fatalf("Failed to load deltas second time: %v", err)
	}

	// Should return the same cached slice
	if len(deltas1) == 0 && len(deltas2) == 0 {
		// Both are empty, which is fine for cache consistency
		return
	}

	if len(deltas1) > 0 && &deltas1[0] != &deltas2[0] {
		t.Error("Expected second load to return cached slice")
	}
}

func TestFileStorage_DeltasCacheUpdateOnSave(t *testing.T) {
	storage, _ := setupTestStorage(t)

	// Create and save initial deltas
	initialDeltas := []core.Delta{
		{
			Action:     core.ActionAdd,
			QuestionID: 1,
			CreatedAt:  testTime,
		},
	}

	err := storage.SaveDeltas(initialDeltas)
	if err != nil {
		t.Fatalf("Failed to save initial deltas: %v", err)
	}

	// Load should return cached version
	_, err = storage.LoadDeltas()
	if err != nil {
		t.Fatalf("Failed to load cached deltas: %v", err)
	}

	// Create and save updated deltas
	updatedDeltas := []core.Delta{
		{
			Action:     core.ActionAdd,
			QuestionID: 1,
			CreatedAt:  testTime,
		},
		{
			Action:     core.ActionUpdate,
			QuestionID: 1,
			CreatedAt:  testTime,
		},
	}

	err = storage.SaveDeltas(updatedDeltas)
	if err != nil {
		t.Fatalf("Failed to save updated deltas: %v", err)
	}

	// Load should return updated cached version
	loadedDeltas, err := storage.LoadDeltas()
	if err != nil {
		t.Fatalf("Failed to load updated deltas: %v", err)
	}

	if len(loadedDeltas) != 2 {
		t.Errorf("Expected 2 deltas after update, got %d", len(loadedDeltas))
	}

	// Should return the same slice as the updated deltas
	if &loadedDeltas[0] != &updatedDeltas[0] {
		t.Error("Expected load to return the same slice as the updated deltas")
	}
}

func TestFileStorage_CacheInvalidation(t *testing.T) {
	storage, _ := setupTestStorage(t)

	// Load to populate cache
	_, err := storage.LoadQuestionStore()
	if err != nil {
		t.Fatalf("Failed to load question store: %v", err)
	}

	_, err = storage.LoadDeltas()
	if err != nil {
		t.Fatalf("Failed to load deltas: %v", err)
	}

	// Invalidate cache
	storage.InvalidateCache()

	// Load again should read from file (not cache)
	store, err := storage.LoadQuestionStore()
	if err != nil {
		t.Fatalf("Failed to load question store after invalidation: %v", err)
	}

	deltas, err := storage.LoadDeltas()
	if err != nil {
		t.Fatalf("Failed to load deltas after invalidation: %v", err)
	}

	// Should return empty data (since files are empty)
	if store.MaxID != 0 {
		t.Errorf("Expected MaxID 0 after cache invalidation, got %d", store.MaxID)
	}

	if len(deltas) != 0 {
		t.Errorf("Expected 0 deltas after cache invalidation, got %d", len(deltas))
	}
}

func TestFileStorage_CacheConsistency(t *testing.T) {
	storage, _ := setupTestStorage(t)

	// Create and save a store
	store := &QuestionStore{
		MaxID:     1,
		Questions: map[int]*core.Question{1: createTestQuestion(1, "test")},
		URLIndex:  map[string]int{"test": 1},
		URLTrie:   search.NewTrie(3),
		NoteTrie:  search.NewTrie(3),
	}

	err := storage.SaveQuestionStore(store)
	if err != nil {
		t.Fatalf("Failed to save store: %v", err)
	}

	// Load multiple times - should return same cached instance
	store1, err := storage.LoadQuestionStore()
	if err != nil {
		t.Fatalf("Failed to load store first time: %v", err)
	}

	store2, err := storage.LoadQuestionStore()
	if err != nil {
		t.Fatalf("Failed to load store second time: %v", err)
	}

	store3, err := storage.LoadQuestionStore()
	if err != nil {
		t.Fatalf("Failed to load store third time: %v", err)
	}

	// All should be the same instance
	if store1 != store2 || store2 != store3 {
		t.Error("Expected all loads to return the same cached instance")
	}

	// Modify the loaded store
	store1.MaxID = 999

	// All references should see the change
	if store2.MaxID != 999 || store3.MaxID != 999 {
		t.Error("Expected all cached references to see the same data")
	}
}

func TestFileStorage_CacheWithFileModification(t *testing.T) {
	storage, _ := setupTestStorage(t)

	// Load to populate cache
	_, err := storage.LoadQuestionStore()
	if err != nil {
		t.Fatalf("Failed to load question store: %v", err)
	}

	// Manually modify the file to simulate external changes
	modifiedStore := &QuestionStore{
		MaxID:     999,
		Questions: map[int]*core.Question{999: createTestQuestion(999, "external")},
		URLIndex:  map[string]int{"external": 999},
		URLTrie:   search.NewTrie(3),
		NoteTrie:  search.NewTrie(3),
	}

	// Save the modified store (this will update the cache)
	err = storage.SaveQuestionStore(modifiedStore)
	if err != nil {
		t.Fatalf("Failed to save modified store: %v", err)
	}

	// Load should return the updated cached version
	cachedStore, err := storage.LoadQuestionStore()
	if err != nil {
		t.Fatalf("Failed to load cached store: %v", err)
	}

	if cachedStore.MaxID != 999 {
		t.Errorf("Expected MaxID 999 from cache, got %d", cachedStore.MaxID)
	}

	// Invalidate cache and load again
	storage.InvalidateCache()
	loadedStore, err := storage.LoadQuestionStore()
	if err != nil {
		t.Fatalf("Failed to load store after invalidation: %v", err)
	}

	if loadedStore.MaxID != 999 {
		t.Errorf("Expected MaxID 999 after cache invalidation, got %d", loadedStore.MaxID)
	}
}

func TestFileStorage_DeleteAllData(t *testing.T) {
	storage, _ := setupTestStorage(t)

	// Create and save some data first
	store := &QuestionStore{
		MaxID:     3,
		Questions: map[int]*core.Question{1: createTestQuestion(1, "test1"), 2: createTestQuestion(2, "test2")},
		URLIndex:  map[string]int{"test1": 1, "test2": 2},
		URLTrie:   search.NewTrie(3),
		NoteTrie:  search.NewTrie(3),
	}
	if err := storage.SaveQuestionStore(store); err != nil {
		t.Fatalf("Failed to save question store: %v", err)
	}

	deltas := []core.Delta{
		{Action: core.ActionAdd, QuestionID: 1, CreatedAt: testTime, NewState: createTestQuestion(1, "test1")},
	}
	if err := storage.SaveDeltas(deltas); err != nil {
		t.Fatalf("Failed to save deltas: %v", err)
	}

	// Verify data exists
	loadedStore, err := storage.LoadQuestionStore()
	if err != nil {
		t.Fatalf("Failed to load store: %v", err)
	}
	if len(loadedStore.Questions) != 2 {
		t.Errorf("Expected 2 questions before delete, got %d", len(loadedStore.Questions))
	}

	loadedDeltas, err := storage.LoadDeltas()
	if err != nil {
		t.Fatalf("Failed to load deltas: %v", err)
	}
	if len(loadedDeltas) != 1 {
		t.Errorf("Expected 1 delta before delete, got %d", len(loadedDeltas))
	}

	// Delete all data
	if err := storage.DeleteAllData(); err != nil {
		t.Fatalf("Failed to delete all data: %v", err)
	}

	// Verify cache is cleared and new load returns empty data
	loadedStore, err = storage.LoadQuestionStore()
	if err != nil {
		t.Fatalf("Failed to load store after delete: %v", err)
	}
	if len(loadedStore.Questions) != 0 {
		t.Errorf("Expected 0 questions after delete, got %d", len(loadedStore.Questions))
	}

	loadedDeltas, err = storage.LoadDeltas()
	if err != nil {
		t.Fatalf("Failed to load deltas after delete: %v", err)
	}
	if len(loadedDeltas) != 0 {
		t.Errorf("Expected 0 deltas after delete, got %d", len(loadedDeltas))
	}
}
