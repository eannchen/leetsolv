package usecase

import (
	"testing"
	"time"

	"leetsolv/config"
	"leetsolv/core"
	"leetsolv/internal/clock"
	"leetsolv/storage"
)

// setupTestEnvironment creates a test environment with temporary files
func setupTestEnvironment(t *testing.T) (*config.TestConfig, *QuestionUseCaseImpl) {
	// Create test configuration with temporary files
	testConfig := config.MockEnv(t)

	// Create mock clock
	mockClock := clock.NewClock()

	// Create storage with test files
	storage := storage.NewFileStorage(testConfig.QuestionsFile, testConfig.DeltasFile)

	// Create scheduler
	scheduler := core.NewSM2Scheduler(mockClock)

	// Create use case
	useCase := NewQuestionUseCase(storage, scheduler, mockClock)

	return testConfig, useCase
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

func TestQuestionUseCase_UpsertQuestion(t *testing.T) {
	// Setup test environment with temporary files
	_, useCase := setupTestEnvironment(t)

	// Test data
	url := "https://leetcode.com/problems/two-sum"
	note := "Test question for unit testing"
	familiarity := core.Medium
	importance := core.MediumImportance

	// Test upserting a new question
	question, err := useCase.UpsertQuestion(url, note, familiarity, importance)
	if err != nil {
		t.Fatalf("Failed to upsert question: %v", err)
	}

	// Verify the question was created correctly
	if question.URL != url {
		t.Errorf("Expected URL %s, got %s", url, question.URL)
	}

	if question.Note != note {
		t.Errorf("Expected note %s, got %s", note, question.Note)
	}

	if question.Familiarity != familiarity {
		t.Errorf("Expected familiarity %d, got %d", familiarity, question.Familiarity)
	}

	if question.Importance != importance {
		t.Errorf("Expected importance %d, got %d", importance, question.Importance)
	}

	// Verify the question was saved to the test file
	store, err := useCase.Storage.LoadQuestionStore()
	if err != nil {
		t.Fatalf("Failed to load question store: %v", err)
	}

	if len(store.Questions) != 1 {
		t.Errorf("Expected 1 question in storage, got %d", len(store.Questions))
	}

	// Test updating the same question
	updatedNote := "Updated test question"
	updatedFamiliarity := core.Easy
	updatedImportance := core.HighImportance

	updatedQuestion, err := useCase.UpsertQuestion(url, updatedNote, updatedFamiliarity, updatedImportance)
	if err != nil {
		t.Fatalf("Failed to update question: %v", err)
	}

	// Verify the question was updated correctly
	if updatedQuestion.Note != updatedNote {
		t.Errorf("Expected updated note %s, got %s", updatedNote, updatedQuestion.Note)
	}

	if updatedQuestion.Familiarity != updatedFamiliarity {
		t.Errorf("Expected updated familiarity %d, got %d", updatedFamiliarity, updatedQuestion.Familiarity)
	}

	if updatedQuestion.Importance != updatedImportance {
		t.Errorf("Expected updated importance %d, got %d", updatedImportance, updatedQuestion.Importance)
	}
}

func TestQuestionUseCase_GetQuestion(t *testing.T) {
	_, useCase := setupTestEnvironment(t)

	// Create test data
	testQuestion := createTestQuestion(1, "https://leetcode.com/problems/test")
	questions := map[int]*core.Question{
		1: testQuestion,
	}

	// Write test data to temporary files
	store := &storage.QuestionStore{
		Questions: questions,
		URLIndex:  map[string]int{"https://leetcode.com/problems/test": 1},
		MaxID:     1,
	}

	err := useCase.Storage.SaveQuestionStore(store)
	if err != nil {
		t.Fatalf("Failed to save test data: %v", err)
	}

	// Test getting question by ID
	question, err := useCase.GetQuestion("1")
	if err != nil {
		t.Fatalf("Failed to get question by ID: %v", err)
	}

	if question.ID != 1 {
		t.Errorf("Expected question ID 1, got %d", question.ID)
	}

	// Test getting question by URL
	question, err = useCase.GetQuestion("https://leetcode.com/problems/test")
	if err != nil {
		t.Fatalf("Failed to get question by URL: %v", err)
	}

	if question.URL != "https://leetcode.com/problems/test" {
		t.Errorf("Expected URL %s, got %s", "https://leetcode.com/problems/test", question.URL)
	}
}

func TestQuestionUseCase_DeleteQuestion(t *testing.T) {
	_, useCase := setupTestEnvironment(t)

	// Create test data
	testQuestion := createTestQuestion(1, "https://leetcode.com/problems/test")
	questions := map[int]*core.Question{
		1: testQuestion,
	}

	// Write test data to temporary files
	store := &storage.QuestionStore{
		Questions: questions,
		URLIndex:  map[string]int{"https://leetcode.com/problems/test": 1},
		MaxID:     1,
	}

	err := useCase.Storage.SaveQuestionStore(store)
	if err != nil {
		t.Fatalf("Failed to save test data: %v", err)
	}

	// Test deleting question by ID
	deletedQuestion, err := useCase.DeleteQuestion("1")
	if err != nil {
		t.Fatalf("Failed to delete question: %v", err)
	}

	if deletedQuestion.ID != 1 {
		t.Errorf("Expected deleted question ID 1, got %d", deletedQuestion.ID)
	}

	// Verify the question was removed from storage
	store, err = useCase.Storage.LoadQuestionStore()
	if err != nil {
		t.Fatalf("Failed to load question store: %v", err)
	}

	if len(store.Questions) != 0 {
		t.Errorf("Expected 0 questions in storage after deletion, got %d", len(store.Questions))
	}
}

func TestQuestionUseCase_ListQuestionsSummary(t *testing.T) {
	_, useCase := setupTestEnvironment(t)

	// Create test questions with different review dates
	now := time.Now()
	tomorrow := now.Add(24 * time.Hour)
	yesterday := now.Add(-24 * time.Hour)

	dueQuestion := createTestQuestion(1, "https://leetcode.com/problems/due")
	dueQuestion.NextReview = yesterday

	upcomingQuestion := createTestQuestion(2, "https://leetcode.com/problems/upcoming")
	upcomingQuestion.NextReview = tomorrow

	questions := map[int]*core.Question{
		1: dueQuestion,
		2: upcomingQuestion,
	}

	// Write test data to temporary files
	store := &storage.QuestionStore{
		Questions: questions,
		URLIndex: map[string]int{
			"https://leetcode.com/problems/due":      1,
			"https://leetcode.com/problems/upcoming": 2,
		},
		MaxID: 2,
	}

	err := useCase.Storage.SaveQuestionStore(store)
	if err != nil {
		t.Fatalf("Failed to save test data: %v", err)
	}

	// Test listing questions summary
	due, upcoming, total, err := useCase.ListQuestionsSummary()
	if err != nil {
		t.Fatalf("Failed to list questions summary: %v", err)
	}

	if len(due) != 1 {
		t.Errorf("Expected 1 due question, got %d", len(due))
	}

	if len(upcoming) != 1 {
		t.Errorf("Expected 1 upcoming question, got %d", len(upcoming))
	}

	if total != 2 {
		t.Errorf("Expected total 2 questions, got %d", total)
	}
}
