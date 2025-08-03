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

	// Create questions with different review dates
	now := time.Now()
	tomorrow := now.Add(24 * time.Hour)
	yesterday := now.Add(-24 * time.Hour)

	// Create a due question (review date in the past)
	dueURL := "https://leetcode.com/problems/due"
	dueQuestion, err := useCase.UpsertQuestion(dueURL, "Due question", core.Medium, core.MediumImportance)
	if err != nil {
		t.Fatalf("Failed to create due question: %v", err)
	}

	// Manually set the review date to yesterday
	store, err := useCase.Storage.LoadQuestionStore()
	if err != nil {
		t.Fatalf("Failed to load question store: %v", err)
	}

	store.Questions[dueQuestion.ID].NextReview = yesterday
	err = useCase.Storage.SaveQuestionStore(store)
	if err != nil {
		t.Fatalf("Failed to save updated question store: %v", err)
	}

	// Create an upcoming question (review date tomorrow)
	upcomingURL := "https://leetcode.com/problems/upcoming"
	upcomingQuestion, err := useCase.UpsertQuestion(upcomingURL, "Upcoming question", core.Medium, core.MediumImportance)
	if err != nil {
		t.Fatalf("Failed to create upcoming question: %v", err)
	}

	// Manually set the review date to tomorrow
	store, err = useCase.Storage.LoadQuestionStore()
	if err != nil {
		t.Fatalf("Failed to load question store: %v", err)
	}

	store.Questions[upcomingQuestion.ID].NextReview = tomorrow
	err = useCase.Storage.SaveQuestionStore(store)
	if err != nil {
		t.Fatalf("Failed to save updated question store: %v", err)
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

	// Verify the due question
	if summary.TopDue[0].URL != dueURL {
		t.Errorf("Expected due question URL %s, got %s", dueURL, summary.TopDue[0].URL)
	}

	// Verify the upcoming question
	if summary.TopUpcoming[0].URL != upcomingURL {
		t.Errorf("Expected upcoming question URL %s, got %s", upcomingURL, summary.TopUpcoming[0].URL)
	}
}

func TestQuestionUseCase_Integration_PaginateQuestions(t *testing.T) {
	useCase, _ := setupIntegrationTest(t)

	// Create multiple questions
	questions := []string{
		"https://leetcode.com/problems/test1",
		"https://leetcode.com/problems/test2",
		"https://leetcode.com/problems/test3",
		"https://leetcode.com/problems/test4",
		"https://leetcode.com/problems/test5",
	}

	for i, url := range questions {
		_, err := useCase.UpsertQuestion(url, "Test question", core.Medium, core.MediumImportance)
		if err != nil {
			t.Fatalf("Failed to create question %d: %v", i+1, err)
		}
	}

	// Get all questions
	allQuestions, err := useCase.ListQuestionsOrderByDesc()
	if err != nil {
		t.Fatalf("Failed to list questions: %v", err)
	}

	if len(allQuestions) != 5 {
		t.Errorf("Expected 5 questions, got %d", len(allQuestions))
	}

	// Test pagination with page size 2
	pageSize := 2
	page := 0

	// First page
	paginatedQuestions, totalPages, err := useCase.PaginateQuestions(allQuestions, pageSize, page)
	if err != nil {
		t.Fatalf("Failed to paginate questions: %v", err)
	}

	if len(paginatedQuestions) != 2 {
		t.Errorf("Expected 2 questions on first page, got %d", len(paginatedQuestions))
	}

	if totalPages != 3 {
		t.Errorf("Expected 3 total pages, got %d", totalPages)
	}

	// Second page
	page = 1
	paginatedQuestions, totalPages, err = useCase.PaginateQuestions(allQuestions, pageSize, page)
	if err != nil {
		t.Fatalf("Failed to paginate questions: %v", err)
	}

	if len(paginatedQuestions) != 2 {
		t.Errorf("Expected 2 questions on second page, got %d", len(paginatedQuestions))
	}

	// Third page
	page = 2
	paginatedQuestions, totalPages, err = useCase.PaginateQuestions(allQuestions, pageSize, page)
	if err != nil {
		t.Fatalf("Failed to paginate questions: %v", err)
	}

	if len(paginatedQuestions) != 1 {
		t.Errorf("Expected 1 question on third page, got %d", len(paginatedQuestions))
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

	// Test pagination with invalid page
	questions := []core.Question{
		{ID: 1, URL: "https://leetcode.com/problems/test"},
	}
	_, _, err = useCase.PaginateQuestions(questions, 5, -1)
	if err == nil {
		t.Error("Expected error when paginating with invalid page")
	}
}
