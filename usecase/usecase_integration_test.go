package usecase

import (
	"testing"
	"time"

	"leetsolv/config"
	"leetsolv/core"
	"leetsolv/internal/clock"
	"leetsolv/storage"
)

// setupIntegrationTest creates a complete test environment with real dependencies
func setupIntegrationTest(t *testing.T) (*QuestionUseCaseImpl, *config.TestConfig) {
	testConfig := config.MockEnv(t)
	mockClock := clock.NewClock()
	storage := storage.NewFileStorage(testConfig.QuestionsFile, testConfig.DeltasFile)
	scheduler := core.NewSM2Scheduler(mockClock)
	useCase := NewQuestionUseCase(storage, scheduler, mockClock)
	return useCase, testConfig
}

func TestQuestionUseCase_Integration_UpsertAndGet(t *testing.T) {
	useCase, _ := setupIntegrationTest(t)

	// Test upserting a new question
	url := "https://leetcode.com/problems/two-sum"
	note := "Test integration question"
	familiarity := core.Medium
	importance := core.MediumImportance

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

	// Test getting the question by ID
	retrievedQuestion, err := useCase.GetQuestion("1")
	if err != nil {
		t.Fatalf("Failed to get question by ID: %v", err)
	}

	if retrievedQuestion.ID != question.ID {
		t.Errorf("Expected question ID %d, got %d", question.ID, retrievedQuestion.ID)
	}

	// Test getting the question by URL
	retrievedQuestion, err = useCase.GetQuestion(url)
	if err != nil {
		t.Fatalf("Failed to get question by URL: %v", err)
	}

	if retrievedQuestion.URL != url {
		t.Errorf("Expected question URL %s, got %s", url, retrievedQuestion.URL)
	}
}

func TestQuestionUseCase_Integration_UpdateQuestion(t *testing.T) {
	useCase, _ := setupIntegrationTest(t)

	// Create initial question
	url := "https://leetcode.com/problems/test"
	note := "Initial note"
	familiarity := core.Medium
	importance := core.MediumImportance

	question, err := useCase.UpsertQuestion(url, note, familiarity, importance)
	if err != nil {
		t.Fatalf("Failed to create initial question: %v", err)
	}

	// Update the question
	updatedNote := "Updated note"
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

	// Verify the ID remains the same
	if updatedQuestion.ID != question.ID {
		t.Errorf("Expected same ID %d, got %d", question.ID, updatedQuestion.ID)
	}
}

func TestQuestionUseCase_Integration_DeleteAndUndo(t *testing.T) {
	useCase, _ := setupIntegrationTest(t)

	// Create a question
	url := "https://leetcode.com/problems/test"
	note := "Test question for deletion"
	familiarity := core.Medium
	importance := core.MediumImportance

	question, err := useCase.UpsertQuestion(url, note, familiarity, importance)
	if err != nil {
		t.Fatalf("Failed to create question: %v", err)
	}

	// Delete the question
	deletedQuestion, err := useCase.DeleteQuestion("1")
	if err != nil {
		t.Fatalf("Failed to delete question: %v", err)
	}

	if deletedQuestion.ID != question.ID {
		t.Errorf("Expected deleted question ID %d, got %d", question.ID, deletedQuestion.ID)
	}

	// Verify the question is no longer retrievable
	_, err = useCase.GetQuestion("1")
	if err == nil {
		t.Error("Expected error when getting deleted question")
	}

	// Undo the deletion
	err = useCase.Undo()
	if err != nil {
		t.Fatalf("Failed to undo deletion: %v", err)
	}

	// Verify the question is retrievable again
	retrievedQuestion, err := useCase.GetQuestion("1")
	if err != nil {
		t.Fatalf("Failed to get question after undo: %v", err)
	}

	if retrievedQuestion.ID != question.ID {
		t.Errorf("Expected retrieved question ID %d, got %d", question.ID, retrievedQuestion.ID)
	}
}

func TestQuestionUseCase_Integration_ListQuestionsSummary(t *testing.T) {
	useCase, _ := setupIntegrationTest(t)

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

func TestQuestionUseCase_Integration_ErrorHandling(t *testing.T) {
	useCase, _ := setupIntegrationTest(t)

	// Test getting non-existent question
	_, err := useCase.GetQuestion("999")
	if err == nil {
		t.Error("Expected error when getting non-existent question")
	}

	// Test deleting non-existent question
	_, err = useCase.DeleteQuestion("999")
	if err == nil {
		t.Error("Expected error when deleting non-existent question")
	}

	// Test undo when no actions available
	err = useCase.Undo()
	if err == nil {
		t.Error("Expected error when undoing with no actions")
	}
}
