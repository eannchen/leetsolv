package usecase

import (
	"fmt"
	"testing"
	"time"

	"leetsolv/config"
	"leetsolv/core"
	"leetsolv/internal/clock"
	"leetsolv/internal/errs"
	"leetsolv/internal/search"
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
	memory := core.MemoryReasoned

	// Test upserting a new question
	delta, err := useCase.UpsertQuestion(url, note, familiarity, importance, memory)
	if err != nil {
		t.Fatalf("Failed to upsert question: %v", err)
	}

	// Verify the question was created correctly
	if delta.NewState.URL != url {
		t.Errorf("Expected URL %s, got %s", url, delta.NewState.URL)
	}

	if delta.NewState.Note != note {
		t.Errorf("Expected note %s, got %s", note, delta.NewState.Note)
	}

	if delta.NewState.Familiarity != familiarity {
		t.Errorf("Expected familiarity %d, got %d", familiarity, delta.NewState.Familiarity)
	}

	if delta.NewState.Importance != importance {
		t.Errorf("Expected importance %d, got %d", importance, delta.NewState.Importance)
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
	updatedMemory := core.MemoryPartial

	updatedDelta, err := useCase.UpsertQuestion(url, updatedNote, updatedFamiliarity, updatedImportance, updatedMemory)
	if err != nil {
		t.Fatalf("Failed to update question: %v", err)
	}

	// Verify the question was updated correctly
	if updatedDelta.NewState.Note != updatedNote {
		t.Errorf("Expected updated note %s, got %s", updatedNote, updatedDelta.NewState.Note)
	}

	if updatedDelta.NewState.Familiarity != updatedFamiliarity {
		t.Errorf("Expected updated familiarity %d, got %d", updatedFamiliarity, updatedDelta.NewState.Familiarity)
	}

	if updatedDelta.NewState.Importance != updatedImportance {
		t.Errorf("Expected updated importance %d, got %d", updatedImportance, updatedDelta.NewState.Importance)
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
		URLTrie:   search.NewTrie(3),
		NoteTrie:  search.NewTrie(3),
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
		URLTrie:   search.NewTrie(3),
		NoteTrie:  search.NewTrie(3),
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
	tomorrow := now.AddDate(0, 0, 1)   // Add 1 calendar day
	yesterday := now.AddDate(0, 0, -1) // Subtract 1 calendar day

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
		MaxID:    2,
		URLTrie:  search.NewTrie(3),
		NoteTrie: search.NewTrie(3),
	}

	err := useCase.Storage.SaveQuestionStore(store)
	if err != nil {
		t.Fatalf("Failed to save test data: %v", err)
	}

	// Test listing questions summary
	summary, err := useCase.ListQuestionsSummary()
	if err != nil {
		t.Fatalf("Failed to list questions summary: %v", err)
	}

	if len(summary.TopDue) != 1 {
		t.Errorf("Expected 1 due question, got %d", len(summary.TopDue))
	}

	if len(summary.TopUpcoming) != 1 {
		t.Errorf("Expected 1 upcoming question, got %d", len(summary.TopUpcoming))
	}

	if summary.Total != 2 {
		t.Errorf("Expected total 2 questions, got %d", summary.Total)
	}

	if summary.TotalDue != 1 {
		t.Errorf("Expected 1 total due question, got %d", summary.TotalDue)
	}

	if summary.TotalUpcoming != 1 {
		t.Errorf("Expected 1 total upcoming question, got %d", summary.TotalUpcoming)
	}
}

// NEW TESTS FOR BETTER BUG DETECTION

func TestQuestionUseCase_SearchQuestions_WithQueries(t *testing.T) {
	_, useCase := setupTestEnvironment(t)

	// Add questions using the proper method to populate tries
	_, err := useCase.UpsertQuestion("https://leetcode.com/problems/two-sum", "Find two numbers that add up to target", core.Medium, core.MediumImportance, core.MemoryReasoned)
	if err != nil {
		t.Fatalf("Failed to add first question: %v", err)
	}

	_, err = useCase.UpsertQuestion("https://leetcode.com/problems/add-two-numbers", "Add two linked lists representing numbers", core.Medium, core.MediumImportance, core.MemoryReasoned)
	if err != nil {
		t.Fatalf("Failed to add second question: %v", err)
	}

	// Test search with query
	results, err := useCase.SearchQuestions([]string{"two"}, nil)
	if err != nil {
		t.Fatalf("Failed to search questions: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 search results, got %d", len(results))
	}

	// Verify both questions are found (both contain "two")
	foundTwoSum := false
	foundAddTwo := false
	for _, result := range results {
		if result.URL == "https://leetcode.com/problems/two-sum" {
			foundTwoSum = true
		}
		if result.URL == "https://leetcode.com/problems/add-two-numbers" {
			foundAddTwo = true
		}
	}

	if !foundTwoSum {
		t.Error("Expected to find two-sum question")
	}
	if !foundAddTwo {
		t.Error("Expected to find add-two-numbers question")
	}
}

func TestQuestionUseCase_SearchQuestions_WithFilters(t *testing.T) {
	_, useCase := setupTestEnvironment(t)

	// Add questions using the proper method to populate tries
	_, err := useCase.UpsertQuestion("https://leetcode.com/problems/test1", "Test question 1", core.Easy, core.HighImportance, core.MemoryReasoned)
	if err != nil {
		t.Fatalf("Failed to add first question: %v", err)
	}

	_, err = useCase.UpsertQuestion("https://leetcode.com/problems/test2", "Test question 2", core.Hard, core.LowImportance, core.MemoryReasoned)
	if err != nil {
		t.Fatalf("Failed to add second question: %v", err)
	}

	// Test search with familiarity filter
	familiarity := core.Easy
	filter := &core.SearchFilter{Familiarity: &familiarity}
	results, err := useCase.SearchQuestions([]string{}, filter)
	if err != nil {
		t.Fatalf("Failed to search questions with filter: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 search result with familiarity filter, got %d", len(results))
	}

	if results[0].Familiarity != core.Easy {
		t.Errorf("Expected familiarity Easy, got %d", results[0].Familiarity)
	}
}

func TestQuestionUseCase_SearchQuestions_EmptyQuery(t *testing.T) {
	_, useCase := setupTestEnvironment(t)

	// Add questions using the proper method
	_, err := useCase.UpsertQuestion("https://leetcode.com/problems/test1", "Test question 1", core.Medium, core.MediumImportance, core.MemoryReasoned)
	if err != nil {
		t.Fatalf("Failed to add first question: %v", err)
	}

	_, err = useCase.UpsertQuestion("https://leetcode.com/problems/test2", "Test question 2", core.Medium, core.MediumImportance, core.MemoryReasoned)
	if err != nil {
		t.Fatalf("Failed to add second question: %v", err)
	}

	// Test search with empty query (should return all questions)
	results, err := useCase.SearchQuestions([]string{}, nil)
	if err != nil {
		t.Fatalf("Failed to search questions: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 search results for empty query, got %d", len(results))
	}
}

func TestQuestionUseCase_Undo_AddAction(t *testing.T) {
	_, useCase := setupTestEnvironment(t)

	// Add a question
	url := "https://leetcode.com/problems/test"
	note := "Test question for undo"
	familiarity := core.Medium
	importance := core.MediumImportance
	memory := core.MemoryReasoned

	_, err := useCase.UpsertQuestion(url, note, familiarity, importance, memory)
	if err != nil {
		t.Fatalf("Failed to add question: %v", err)
	}

	// Verify question exists
	_, err = useCase.GetQuestion("1")
	if err != nil {
		t.Fatalf("Failed to get question after add: %v", err)
	}

	// Undo the add action
	err = useCase.Undo()
	if err != nil {
		t.Fatalf("Failed to undo add action: %v", err)
	}

	// Verify question no longer exists
	_, err = useCase.GetQuestion("1")
	if err == nil {
		t.Error("Expected error when getting question after undo")
	}
}

func TestQuestionUseCase_Undo_UpdateAction(t *testing.T) {
	_, useCase := setupTestEnvironment(t)

	// Add a question
	url := "https://leetcode.com/problems/test"
	note := "Original note"
	familiarity := core.Medium
	importance := core.MediumImportance
	memory := core.MemoryReasoned

	_, err := useCase.UpsertQuestion(url, note, familiarity, importance, memory)
	if err != nil {
		t.Fatalf("Failed to add question: %v", err)
	}

	// Update the question
	updatedNote := "Updated note"
	updatedFamiliarity := core.Easy
	updatedImportance := core.HighImportance
	updatedMemory := core.MemoryPartial

	_, err = useCase.UpsertQuestion(url, updatedNote, updatedFamiliarity, updatedImportance, updatedMemory)
	if err != nil {
		t.Fatalf("Failed to update question: %v", err)
	}

	// Verify the update
	retrievedQuestion, err := useCase.GetQuestion("1")
	if err != nil {
		t.Fatalf("Failed to get updated question: %v", err)
	}

	if retrievedQuestion.Note != updatedNote {
		t.Errorf("Expected updated note %s, got %s", updatedNote, retrievedQuestion.Note)
	}

	// Undo the update action
	err = useCase.Undo()
	if err != nil {
		t.Fatalf("Failed to undo update action: %v", err)
	}

	// Verify the original state is restored
	retrievedQuestion, err = useCase.GetQuestion("1")
	if err != nil {
		t.Fatalf("Failed to get question after undo: %v", err)
	}

	if retrievedQuestion.Note != note {
		t.Errorf("Expected original note %s, got %s", note, retrievedQuestion.Note)
	}
}

func TestQuestionUseCase_Undo_DeleteAction(t *testing.T) {
	_, useCase := setupTestEnvironment(t)

	// Add a question
	url := "https://leetcode.com/problems/test"
	note := "Test question for delete undo"
	familiarity := core.Medium
	importance := core.MediumImportance
	memory := core.MemoryReasoned

	delta, err := useCase.UpsertQuestion(url, note, familiarity, importance, memory)
	if err != nil {
		t.Fatalf("Failed to add question: %v", err)
	}

	// Delete the question
	deletedQuestion, err := useCase.DeleteQuestion("1")
	if err != nil {
		t.Fatalf("Failed to delete question: %v", err)
	}

	if deletedQuestion.ID != delta.NewState.ID {
		t.Errorf("Expected deleted question ID %d, got %d", delta.NewState.ID, deletedQuestion.ID)
	}

	// Verify question no longer exists
	_, err = useCase.GetQuestion("1")
	if err == nil {
		t.Error("Expected error when getting deleted question")
	}

	// Undo the delete action
	err = useCase.Undo()
	if err != nil {
		t.Fatalf("Failed to undo delete action: %v", err)
	}

	// Verify question exists again
	retrievedQuestion, err := useCase.GetQuestion("1")
	if err != nil {
		t.Fatalf("Failed to get question after undo: %v", err)
	}

	if retrievedQuestion.ID != delta.NewState.ID {
		t.Errorf("Expected restored question ID %d, got %d", delta.NewState.ID, retrievedQuestion.ID)
	}
}

func TestQuestionUseCase_Undo_NoActions(t *testing.T) {
	_, useCase := setupTestEnvironment(t)

	// Try to undo when no actions are available
	err := useCase.Undo()
	if err == nil {
		t.Error("Expected error when undoing with no actions")
	}
}

func TestQuestionUseCase_GetQuestion_NotFound(t *testing.T) {
	_, useCase := setupTestEnvironment(t)

	// Try to get non-existent question
	_, err := useCase.GetQuestion("999")
	if err == nil {
		t.Error("Expected error when getting non-existent question")
	}
}

func TestQuestionUseCase_DeleteQuestion_NotFound(t *testing.T) {
	_, useCase := setupTestEnvironment(t)

	// Try to delete non-existent question
	_, err := useCase.DeleteQuestion("999")
	if err == nil {
		t.Error("Expected error when deleting non-existent question")
	}
}

func TestQuestionUseCase_InvalidURL(t *testing.T) {
	_, useCase := setupTestEnvironment(t)

	// Try to upsert with invalid URL
	_, err := useCase.UpsertQuestion("invalid-url", "test", core.Medium, core.MediumImportance, core.MemoryReasoned)
	if err == nil {
		t.Error("Expected error when upserting with invalid URL")
	}
}

func TestQuestionUseCase_ConcurrentAccess(t *testing.T) {
	_, useCase := setupTestEnvironment(t)

	// Test concurrent access to use case
	done := make(chan bool, 3)

	// Goroutine 1: Add questions
	go func() {
		defer func() { done <- true }()
		for i := 0; i < 5; i++ {
			url := "https://leetcode.com/problems/test"
			_, err := useCase.UpsertQuestion(url, "test", core.Medium, core.MediumImportance, core.MemoryReasoned)
			if err != nil {
				t.Errorf("Failed to add question in goroutine 1: %v", err)
			}
		}
	}()

	// Goroutine 2: Get questions (with proper error handling)
	go func() {
		defer func() { done <- true }()
		for i := 0; i < 5; i++ {
			_, err := useCase.GetQuestion("1")
			if err != nil && err != errs.ErrQuestionNotFound && err != errs.ErrNoQuestionsAvailable {
				t.Errorf("Failed to get question in goroutine 2: %v", err)
			}
		}
	}()

	// Goroutine 3: List summary
	go func() {
		defer func() { done <- true }()
		for i := 0; i < 5; i++ {
			_, err := useCase.ListQuestionsSummary()
			if err != nil {
				t.Errorf("Failed to list summary in goroutine 3: %v", err)
			}
		}
	}()

	// Wait for all goroutines to complete
	<-done
	<-done
	<-done
}

func TestQuestionUseCase_LargeDataSet(t *testing.T) {
	_, useCase := setupTestEnvironment(t)

	// Add many questions with unique URLs
	for i := 0; i < 10; i++ {
		url := fmt.Sprintf("https://leetcode.com/problems/test%d", i)
		_, err := useCase.UpsertQuestion(url, "test", core.Medium, core.MediumImportance, core.MemoryReasoned)
		if err != nil {
			t.Fatalf("Failed to add question %d: %v", i, err)
		}
	}

	// Test listing all questions
	questions, err := useCase.ListQuestionsOrderByDesc()
	if err != nil {
		t.Fatalf("Failed to list questions: %v", err)
	}

	if len(questions) != 10 {
		t.Errorf("Expected 10 questions, got %d", len(questions))
	}

	// Test search with large dataset
	results, err := useCase.SearchQuestions([]string{"test"}, nil)
	if err != nil {
		t.Fatalf("Failed to search large dataset: %v", err)
	}

	if len(results) != 10 {
		t.Errorf("Expected 10 search results, got %d", len(results))
	}
}

func TestQuestionUseCase_SchedulerIntegration(t *testing.T) {
	_, useCase := setupTestEnvironment(t)

	// Add a question with different familiarity levels
	testCases := []struct {
		familiarity core.Familiarity
		expected    string
	}{
		{core.VeryHard, "Expected very hard question to be scheduled"},
		{core.Hard, "Expected hard question to be scheduled"},
		{core.Medium, "Expected medium question to be scheduled"},
		{core.Easy, "Expected easy question to be scheduled"},
		{core.VeryEasy, "Expected very easy question to be scheduled"},
	}

	for i, tc := range testCases {
		url := fmt.Sprintf("https://leetcode.com/problems/test%d", i)
		delta, err := useCase.UpsertQuestion(url, "test", tc.familiarity, core.MediumImportance, core.MemoryReasoned)
		if err != nil {
			t.Fatalf("Failed to add question with familiarity %d: %v", tc.familiarity, err)
		}

		// Verify the question was scheduled
		if delta.NewState.NextReview.Before(time.Now()) {
			t.Errorf("Expected question to be scheduled in the future, got %v", delta.NewState.NextReview)
		}

		// Clean up for next test
		_, err = useCase.DeleteQuestion(fmt.Sprintf("%d", delta.NewState.ID))
		if err != nil {
			t.Fatalf("Failed to delete question: %v", err)
		}
	}
}
