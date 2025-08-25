package handler

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/eannchen/leetsolv/config"
	"github.com/eannchen/leetsolv/core"
	"github.com/eannchen/leetsolv/internal/errs"
	"github.com/eannchen/leetsolv/internal/logger"
	"github.com/eannchen/leetsolv/usecase"
)

// MockIOHandler implements IOHandler for testing
type MockIOHandler struct {
	output     *bytes.Buffer
	input      *strings.Reader
	readCalls  []string
	writeCalls []string
	lines      []string
	lineIndex  int
}

func NewMockIOHandler(input string) *MockIOHandler {
	lines := strings.Split(input, "\n")
	return &MockIOHandler{
		output: &bytes.Buffer{},
		input:  strings.NewReader(input),
		lines:  lines,
	}
}

func (m *MockIOHandler) Println(a ...interface{}) {
	m.writeCalls = append(m.writeCalls, "Println")
	m.output.WriteString(fmt.Sprintln(a...))
}

func (m *MockIOHandler) Printf(format string, a ...interface{}) {
	m.writeCalls = append(m.writeCalls, "Printf")
	m.output.WriteString(fmt.Sprintf(format, a...))
}

func (m *MockIOHandler) PrintlnColored(color string, a ...interface{}) {
	m.writeCalls = append(m.writeCalls, "PrintlnColored")
	m.output.WriteString(fmt.Sprintln(a...))
}

func (m *MockIOHandler) PrintfColored(color string, format string, a ...interface{}) {
	m.writeCalls = append(m.writeCalls, "PrintfColored")
	m.output.WriteString(fmt.Sprintf(format, a...))
}

func (m *MockIOHandler) ReadLine(scanner *bufio.Scanner, prompt string) string {
	m.readCalls = append(m.readCalls, prompt)
	if m.lineIndex < len(m.lines) {
		line := m.lines[m.lineIndex]
		m.lineIndex++
		return line
	}
	return ""
}

func (m *MockIOHandler) PrintQuestionBrief(q *core.Question) {
	m.writeCalls = append(m.writeCalls, "PrintQuestionBrief")
	m.output.WriteString(fmt.Sprintf("ID: %d, URL: %s\n", q.ID, q.URL))
}

func (m *MockIOHandler) PrintQuestionDetail(question *core.Question) {
	m.writeCalls = append(m.writeCalls, "PrintQuestionDetail")
	m.output.WriteString(fmt.Sprintf("Question Detail - ID: %d, URL: %s\n", question.ID, question.URL))
}

func (m *MockIOHandler) PrintQuestionUpsertDetail(delta *core.Delta) {
	m.writeCalls = append(m.writeCalls, "PrintQuestionUpsertDetail")
	if delta != nil && delta.NewState != nil {
		m.output.WriteString(fmt.Sprintf("Upserted - ID: %d, URL: %s\n", delta.NewState.ID, delta.NewState.URL))
	} else {
		m.output.WriteString("Upserted - <nil>\n")
	}
}

func (m *MockIOHandler) PrintSuccess(message string) {
	m.writeCalls = append(m.writeCalls, "PrintSuccess")
	m.output.WriteString(fmt.Sprintf("SUCCESS: %s\n", message))
}

func (m *MockIOHandler) PrintError(err error) {
	m.writeCalls = append(m.writeCalls, "PrintError")
	m.output.WriteString(fmt.Sprintf("ERROR: %v\n", err))
}

func (m *MockIOHandler) PrintCancel(message string) {
	m.writeCalls = append(m.writeCalls, "PrintCancel")
	m.output.WriteString(fmt.Sprintf("CANCELLED: %s\n", message))
}

// MockQuestionUseCase implements QuestionUseCase for testing
type MockQuestionUseCase struct {
	questions     []core.Question
	shouldError   bool
	errorToReturn error
	upserted      *core.Delta
	deleted       *core.Question
	summary       usecase.QuestionsSummary
	searchResults []core.Question
	pagination    map[string]interface{} // For testing pagination edge cases
}

func NewMockQuestionUseCase() *MockQuestionUseCase {
	return &MockQuestionUseCase{
		questions: []core.Question{},
		summary: usecase.QuestionsSummary{
			TopDue:        []core.Question{},
			TotalDue:      0,
			TopUpcoming:   []core.Question{},
			TotalUpcoming: 0,
			Total:         0,
		},
		pagination: make(map[string]interface{}),
	}
}

func (m *MockQuestionUseCase) ListQuestionsSummary() (usecase.QuestionsSummary, error) {
	if m.shouldError {
		return usecase.QuestionsSummary{}, m.errorToReturn
	}
	return m.summary, nil
}

func (m *MockQuestionUseCase) ListQuestionsOrderByDesc() ([]core.Question, error) {
	if m.shouldError {
		return nil, m.errorToReturn
	}
	return m.questions, nil
}

func (m *MockQuestionUseCase) PaginateQuestions(questions []core.Question, pageSize, page int) ([]core.Question, int, error) {
	if m.shouldError {
		return nil, 0, m.errorToReturn
	}

	// Test edge cases for pagination
	if page < 0 {
		return nil, 0, errs.ErrInvalidPageNumber
	}

	if len(questions) == 0 {
		return nil, 0, nil
	}

	totalPages := (len(questions) + pageSize - 1) / pageSize
	if page >= totalPages {
		return nil, totalPages, errs.ErrInvalidPageNumber
	}

	start := page * pageSize
	end := start + pageSize
	if end > len(questions) {
		end = len(questions)
	}

	return questions[start:end], totalPages, nil
}

func (m *MockQuestionUseCase) GetQuestion(target string) (*core.Question, error) {
	if m.shouldError {
		return nil, m.errorToReturn
	}
	if len(m.questions) > 0 {
		return &m.questions[0], nil
	}
	return nil, errs.ErrQuestionNotFound
}

func (m *MockQuestionUseCase) SearchQuestions(queries []string, filter *core.SearchFilter) ([]core.Question, error) {
	if m.shouldError {
		return nil, m.errorToReturn
	}
	return m.searchResults, nil
}

func (m *MockQuestionUseCase) UpsertQuestion(url, note string, familiarity core.Familiarity, importance core.Importance, memory core.MemoryUse) (*core.Delta, error) {
	if m.shouldError {
		return nil, m.errorToReturn
	}
	return m.upserted, nil
}

func (m *MockQuestionUseCase) DeleteQuestion(target string) (*core.Question, error) {
	if m.shouldError {
		return nil, m.errorToReturn
	}
	return m.deleted, nil
}

func (m *MockQuestionUseCase) Undo() error {
	if m.shouldError {
		return m.errorToReturn
	}
	return nil
}

func (m *MockQuestionUseCase) GetHistory() ([]core.Delta, error) {
	if m.shouldError {
		return nil, m.errorToReturn
	}
	// Return a sample delta for testing
	return []core.Delta{
		{
			Action:     core.ActionAdd,
			QuestionID: 1,
			NewState: &core.Question{
				ID:  1,
				URL: "https://leetcode.com/problems/test-question/",
			},
			CreatedAt: time.Now(),
		},
	}, nil
}

func (m *MockQuestionUseCase) GetSettings() error {
	if m.shouldError {
		return m.errorToReturn
	}
	return nil
}

func (m *MockQuestionUseCase) UpdateSetting(settingName string, value interface{}) error {
	if m.shouldError {
		return m.errorToReturn
	}
	return nil
}

// setupTestHandler creates a test handler with mocked dependencies
func setupTestHandler(t *testing.T) (*HandlerImpl, *MockIOHandler, *MockQuestionUseCase) {
	_, cfg := config.MockEnv(t)
	logger := logger.NewLogger(cfg.InfoLogFile, cfg.ErrorLogFile)
	mockIO := NewMockIOHandler("")
	mockUseCase := NewMockQuestionUseCase()
	handler := NewHandler(cfg, logger, mockUseCase, mockIO, "test-version")
	return handler, mockIO, mockUseCase
}

func TestHandler_HandleList_EmptyQuestions(t *testing.T) {
	handler, mockIO, mockUseCase := setupTestHandler(t)

	// Set up empty questions
	mockUseCase.questions = []core.Question{}

	scanner := bufio.NewScanner(strings.NewReader(""))
	handler.HandleList(scanner)

	// Verify that error was printed
	found := false
	for _, call := range mockIO.writeCalls {
		if call == "PrintError" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected PrintError to be called for empty questions")
	}
}

func TestHandler_HandleList_WithQuestions(t *testing.T) {
	handler, mockIO, mockUseCase := setupTestHandler(t)

	// Set up test questions
	mockUseCase.questions = []core.Question{
		{
			ID:         1,
			URL:        "https://leetcode.com/problems/test1",
			Note:       "Test question 1",
			NextReview: time.Now(),
		},
		{
			ID:         2,
			URL:        "https://leetcode.com/problems/test2",
			Note:       "Test question 2",
			NextReview: time.Now(),
		},
	}

	scanner := bufio.NewScanner(strings.NewReader("q\n"))
	handler.HandleList(scanner)

	// Verify that questions were displayed
	found := false
	for _, call := range mockIO.writeCalls {
		if call == "PrintfColored" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected PrintfColored to be called for displaying questions")
	}
}

func TestHandler_HandleSearch_WithQueries(t *testing.T) {
	handler, mockIO, mockUseCase := setupTestHandler(t)

	// Set up search results
	mockUseCase.searchResults = []core.Question{
		{
			ID:   1,
			URL:  "https://leetcode.com/problems/test1",
			Note: "Test question 1",
		},
	}

	scanner := bufio.NewScanner(strings.NewReader("q\n"))
	handler.HandleSearch(scanner, []string{"test"})

	// Verify that search was performed
	found := false
	for _, call := range mockIO.writeCalls {
		if call == "PrintQuestionBrief" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected PrintQuestionBrief to be called for search results")
	}
}

func TestHandler_HandleSearch_WithFilters(t *testing.T) {
	handler, mockIO, mockUseCase := setupTestHandler(t)

	// Set up search results
	mockUseCase.searchResults = []core.Question{
		{
			ID:   1,
			URL:  "https://leetcode.com/problems/test1",
			Note: "Test question 1",
		},
	}

	scanner := bufio.NewScanner(strings.NewReader("q\n"))
	handler.HandleSearch(scanner, []string{"--familiarity=3", "--importance=2"})

	// Verify that search was performed
	found := false
	for _, call := range mockIO.writeCalls {
		if call == "PrintQuestionBrief" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected PrintQuestionBrief to be called for search results")
	}
}

func TestHandler_HandleSearch_InvalidFilter(t *testing.T) {
	handler, mockIO, _ := setupTestHandler(t)

	scanner := bufio.NewScanner(strings.NewReader(""))
	handler.HandleSearch(scanner, []string{"--familiarity=invalid"})

	// Verify that error was printed
	found := false
	for _, call := range mockIO.writeCalls {
		if call == "PrintError" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected PrintError to be called for invalid filter")
	}
}

func TestHandler_HandleGet_Success(t *testing.T) {
	handler, mockIO, mockUseCase := setupTestHandler(t)

	// Set up test question
	testQuestion := core.Question{
		ID:         1,
		URL:        "https://leetcode.com/problems/test",
		Note:       "Test question",
		NextReview: time.Now(),
	}
	mockUseCase.questions = []core.Question{testQuestion}

	scanner := bufio.NewScanner(strings.NewReader(""))
	handler.HandleGet(scanner, "1")

	// Verify that question details were printed
	found := false
	for _, call := range mockIO.writeCalls {
		if call == "PrintQuestionDetail" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected PrintQuestionDetail to be called")
	}
}

func TestHandler_HandleGet_Error(t *testing.T) {
	handler, mockIO, mockUseCase := setupTestHandler(t)

	// Set up error
	mockUseCase.shouldError = true
	mockUseCase.errorToReturn = errs.ErrQuestionNotFound

	scanner := bufio.NewScanner(strings.NewReader(""))
	handler.HandleGet(scanner, "999")

	// Verify that error was printed
	found := false
	for _, call := range mockIO.writeCalls {
		if call == "PrintError" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected PrintError to be called for error case")
	}
}

func TestHandler_HandleGet_EmptyInput(t *testing.T) {
	handler, mockIO, _ := setupTestHandler(t)

	scanner := bufio.NewScanner(strings.NewReader("\n"))
	handler.HandleGet(scanner, "")

	// Verify that error was printed for empty input
	found := false
	for _, call := range mockIO.writeCalls {
		if call == "PrintError" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected PrintError to be called for empty input")
	}
}

func TestHandler_HandleStatus_Success(t *testing.T) {
	handler, mockIO, mockUseCase := setupTestHandler(t)

	// Set up test summary
	mockUseCase.summary = usecase.QuestionsSummary{
		TopDue: []core.Question{
			{
				ID:   1,
				URL:  "https://leetcode.com/problems/test1",
				Note: "Test question 1",
			},
		},
		TotalDue: 1,
		TopUpcoming: []core.Question{
			{
				ID:   2,
				URL:  "https://leetcode.com/problems/test2",
				Note: "Test question 2",
			},
		},
		TotalUpcoming: 1,
		Total:         2,
	}

	handler.HandleStatus()

	// Verify that status was displayed
	found := false
	for _, call := range mockIO.writeCalls {
		if call == "PrintlnColored" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected PrintlnColored to be called for status display")
	}
}

func TestHandler_HandleStatus_Error(t *testing.T) {
	handler, mockIO, mockUseCase := setupTestHandler(t)

	// Set up error
	mockUseCase.shouldError = true
	mockUseCase.errorToReturn = errs.ErrQuestionNotFound

	handler.HandleStatus()

	// Verify that error was printed
	found := false
	for _, call := range mockIO.writeCalls {
		if call == "PrintError" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected PrintError to be called for status error")
	}
}

func TestHandler_HandleUpsert_Success(t *testing.T) {
	// Create mock IO with proper input
	mockIO := NewMockIOHandler("https://leetcode.com/problems/two-sum\nTest question\n3\n2\n")
	mockUseCase := NewMockQuestionUseCase()
	_, cfg := config.MockEnv(t)
	logger := logger.NewLogger(cfg.InfoLogFile, cfg.ErrorLogFile)
	handler := NewHandler(cfg, logger, mockUseCase, mockIO, "test-version")

	// Set up successful upsert
	upsertedQuestion := &core.Question{
		ID:          1,
		URL:         "https://leetcode.com/problems/test",
		Note:        "Test question",
		Familiarity: core.Medium,
		Importance:  core.MediumImportance,
	}
	mockUseCase.upserted = &core.Delta{
		Action:     core.ActionAdd,
		QuestionID: upsertedQuestion.ID,
		OldState:   nil,
		NewState:   upsertedQuestion,
		CreatedAt:  time.Now(),
	}

	// Test the URL normalization directly first
	normalizedURL, err := handler.normalizeLeetCodeURL("https://leetcode.com/problems/two-sum")
	if err != nil {
		t.Fatalf("URL normalization failed: %v", err)
	}
	t.Logf("Normalized URL: %s", normalizedURL)

	// Simulate user input: URL, note, familiarity (3), importance (2)
	scanner := bufio.NewScanner(strings.NewReader(""))
	handler.HandleUpsert(scanner, "")

	// Verify that success message was printed
	found := false
	for _, call := range mockIO.writeCalls {
		if call == "PrintSuccess" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected PrintSuccess to be called for success message")
	}

	// Also verify that PrintQuestionUpsertDetail was called
	found = false
	for _, call := range mockIO.writeCalls {
		if call == "PrintQuestionUpsertDetail" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected PrintQuestionUpsertDetail to be called")
	}

	// Verify reassurance and normalized URL outputs
	output := mockIO.output.String()
	if !strings.Contains(output, "Provided URL will be normalized to a canonical form to match existing data.") {
		t.Error("Expected reassurance line about URL normalization to be printed")
	}
	if !strings.Contains(output, "Using normalized URL: https://leetcode.com/problems/two-sum/") {
		t.Error("Expected normalized URL to be printed")
	}
}

func TestHandler_HandleUpsert_InvalidURL(t *testing.T) {
	handler, mockIO, _ := setupTestHandler(t)

	// Simulate invalid URL input
	input := "invalid-url\n"
	scanner := bufio.NewScanner(strings.NewReader(input))
	handler.HandleUpsert(scanner, "")

	// Verify that error was printed
	found := false
	for _, call := range mockIO.writeCalls {
		if call == "PrintError" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected PrintError to be called for invalid URL")
	}
}

func TestHandler_HandleUpsert_InvalidFamiliarity(t *testing.T) {
	handler, mockIO, _ := setupTestHandler(t)

	// Simulate valid URL but invalid familiarity
	input := "https://leetcode.com/problems/test\nTest question\n6\n"
	scanner := bufio.NewScanner(strings.NewReader(input))
	handler.HandleUpsert(scanner, "")

	// Verify that error was printed
	found := false
	for _, call := range mockIO.writeCalls {
		if call == "PrintError" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected PrintError to be called for invalid familiarity")
	}
}

func TestHandler_HandleUpsert_InvalidImportance(t *testing.T) {
	handler, mockIO, _ := setupTestHandler(t)

	// Simulate valid URL and familiarity but invalid importance
	input := "https://leetcode.com/problems/test\nTest question\n3\n5\n"
	scanner := bufio.NewScanner(strings.NewReader(input))
	handler.HandleUpsert(scanner, "")

	// Verify that error was printed
	found := false
	for _, call := range mockIO.writeCalls {
		if call == "PrintError" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected PrintError to be called for invalid importance")
	}
}

func TestHandler_HandleUpsert_UseCaseError(t *testing.T) {
	handler, mockIO, mockUseCase := setupTestHandler(t)

	// Set up use case error
	mockUseCase.shouldError = true
	mockUseCase.errorToReturn = errs.ErrQuestionNotFound

	// Simulate valid input
	input := "https://leetcode.com/problems/test\nTest question\n3\n2\n"
	scanner := bufio.NewScanner(strings.NewReader(input))
	handler.HandleUpsert(scanner, "")

	// Verify that error was printed
	found := false
	for _, call := range mockIO.writeCalls {
		if call == "PrintError" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected PrintError to be called for use case error")
	}
}

func TestHandler_HandleDelete_Success(t *testing.T) {
	// Create mock IO with proper input
	mockIO := NewMockIOHandler("y\n")
	mockUseCase := NewMockQuestionUseCase()
	_, cfg := config.MockEnv(t)
	logger := logger.NewLogger(cfg.InfoLogFile, cfg.ErrorLogFile)
	handler := NewHandler(cfg, logger, mockUseCase, mockIO, "test-version")

	// Set up successful deletion
	deletedQuestion := &core.Question{
		ID:   1,
		URL:  "https://leetcode.com/problems/test",
		Note: "Test question",
	}
	mockUseCase.deleted = deletedQuestion

	// Simulate user confirmation
	scanner := bufio.NewScanner(strings.NewReader(""))
	handler.HandleDelete(scanner, "1")

	// Verify that deletion message was printed
	found := false
	for _, call := range mockIO.writeCalls {
		if call == "PrintSuccess" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected PrintSuccess to be called for deletion message")
	}
}

func TestHandler_HandleDelete_Cancelled(t *testing.T) {
	handler, mockIO, _ := setupTestHandler(t)

	// Simulate user cancellation
	input := "n\n"
	scanner := bufio.NewScanner(strings.NewReader(input))
	handler.HandleDelete(scanner, "1")

	// Verify that cancellation message was printed
	found := false
	for _, call := range mockIO.writeCalls {
		if call == "PrintCancel" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected PrintCancel to be called for cancellation message")
	}
}

func TestHandler_HandleDelete_EmptyInput(t *testing.T) {
	handler, mockIO, _ := setupTestHandler(t)

	scanner := bufio.NewScanner(strings.NewReader("\n"))
	handler.HandleDelete(scanner, "")

	// Verify that error was printed for empty input
	found := false
	for _, call := range mockIO.writeCalls {
		if call == "PrintError" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected PrintError to be called for empty input")
	}
}

func TestHandler_HandleDelete_UseCaseError(t *testing.T) {
	// Create mock IO with proper input
	mockIO := NewMockIOHandler("y\n")
	mockUseCase := NewMockQuestionUseCase()
	_, cfg := config.MockEnv(t)
	logger := logger.NewLogger(cfg.InfoLogFile, cfg.ErrorLogFile)
	handler := NewHandler(cfg, logger, mockUseCase, mockIO, "test-version")

	// Set up use case error
	mockUseCase.shouldError = true
	mockUseCase.errorToReturn = errs.ErrQuestionNotFound

	// Simulate user confirmation
	scanner := bufio.NewScanner(strings.NewReader(""))
	handler.HandleDelete(scanner, "999")

	// Verify that error was printed
	found := false
	for _, call := range mockIO.writeCalls {
		if call == "PrintError" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected PrintError to be called for use case error")
	}
}

func TestHandler_HandleUndo_Success(t *testing.T) {
	// Create mock IO with proper input
	mockIO := NewMockIOHandler("y\n")
	mockUseCase := NewMockQuestionUseCase()
	_, cfg := config.MockEnv(t)
	logger := logger.NewLogger(cfg.InfoLogFile, cfg.ErrorLogFile)
	handler := NewHandler(cfg, logger, mockUseCase, mockIO, "test-version")

	// Simulate user confirmation
	scanner := bufio.NewScanner(strings.NewReader(""))
	handler.HandleUndo(scanner)

	// Verify that success message was printed
	found := false
	for _, call := range mockIO.writeCalls {
		if call == "PrintSuccess" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected PrintSuccess to be called for undo success message")
	}
}

func TestHandler_HandleUndo_Error(t *testing.T) {
	// Create mock IO with proper input
	mockIO := NewMockIOHandler("y\n")
	mockUseCase := NewMockQuestionUseCase()
	_, cfg := config.MockEnv(t)
	logger := logger.NewLogger(cfg.InfoLogFile, cfg.ErrorLogFile)
	handler := NewHandler(cfg, logger, mockUseCase, mockIO, "test-version")

	// Set up error
	mockUseCase.shouldError = true
	mockUseCase.errorToReturn = errs.ErrNoActionsToUndo

	// Simulate user confirmation
	scanner := bufio.NewScanner(strings.NewReader(""))
	handler.HandleUndo(scanner)

	// Verify that error was printed
	found := false
	for _, call := range mockIO.writeCalls {
		if call == "PrintError" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected PrintError to be called for undo error")
	}
}

func TestHandler_HandleUndo_Cancelled(t *testing.T) {
	handler, mockIO, _ := setupTestHandler(t)

	// Simulate user cancellation
	input := "n\n"
	scanner := bufio.NewScanner(strings.NewReader(input))
	handler.HandleUndo(scanner)

	// Verify that cancellation message was printed
	found := false
	for _, call := range mockIO.writeCalls {
		if call == "PrintCancel" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected PrintCancel to be called for undo cancellation")
	}
}

func TestHandler_ValidateFamiliarity(t *testing.T) {
	handler, _, _ := setupTestHandler(t)

	// Test valid inputs
	testCases := []struct {
		input    string
		expected core.Familiarity
		hasError bool
	}{
		{"1", core.VeryHard, false},
		{"2", core.Hard, false},
		{"3", core.Medium, false},
		{"4", core.Easy, false},
		{"5", core.VeryEasy, false},
		{"0", 0, true},
		{"6", 0, true},
		{"abc", 0, true},
		{"", 0, true},
		{"-1", 0, true},
	}

	for _, tc := range testCases {
		result, err := handler.validateFamiliarity(tc.input)
		if tc.hasError && err == nil {
			t.Errorf("Expected error for input %s, got none", tc.input)
		}
		if !tc.hasError && err != nil {
			t.Errorf("Expected no error for input %s, got %v", tc.input, err)
		}
		if !tc.hasError && result != tc.expected {
			t.Errorf("Expected %d for input %s, got %d", tc.expected, tc.input, result)
		}
	}
}

func TestHandler_ValidateImportance(t *testing.T) {
	handler, _, _ := setupTestHandler(t)

	// Test valid inputs
	testCases := []struct {
		input    string
		expected core.Importance
		hasError bool
	}{
		{"1", core.LowImportance, false},
		{"2", core.MediumImportance, false},
		{"3", core.HighImportance, false},
		{"4", core.CriticalImportance, false},
		{"0", 0, true},
		{"5", 0, true},
		{"abc", 0, true},
		{"", 0, true},
		{"-1", 0, true},
	}

	for _, tc := range testCases {
		result, err := handler.validateImportance(tc.input)
		if tc.hasError && err == nil {
			t.Errorf("Expected error for input %s, got none", tc.input)
		}
		if !tc.hasError && err != nil {
			t.Errorf("Expected no error for input %s, got %v", tc.input, err)
		}
		if !tc.hasError && result != tc.expected {
			t.Errorf("Expected %d for input %s, got %d", tc.expected, tc.input, result)
		}
	}
}

func TestHandler_NormalizeLeetCodeURL(t *testing.T) {
	handler, _, _ := setupTestHandler(t)

	testCases := []struct {
		input    string
		expected string
		hasError bool
	}{
		{
			"https://leetcode.com/problems/two-sum/",
			"https://leetcode.com/problems/two-sum/",
			false,
		},
		{
			"https://leetcode.com/problems/two-sum",
			"https://leetcode.com/problems/two-sum/",
			false,
		},
		{
			"https://leetcode.com/problems/two-sum/solution/",
			"https://leetcode.com/problems/two-sum/",
			false,
		},
		{
			"https://leetcode.com/problems/two-sum/discuss/",
			"https://leetcode.com/problems/two-sum/",
			false,
		},
		{
			"https://leetcode.com/problems/",
			"",
			true,
		},
		{
			"https://leetcode.com/",
			"",
			true,
		},
		{
			"https://google.com/problems/test",
			"",
			true,
		},
		{
			"invalid-url",
			"",
			true,
		},
		{
			"",
			"",
			true,
		},
		{
			"https://leetcode.com/problems/",
			"",
			true,
		},
		{
			"https://leetcode.com/problems//",
			"",
			true,
		},
		{
			"https://leetcode.com/problems/",
			"",
			true,
		},
	}

	for _, tc := range testCases {
		result, err := handler.normalizeLeetCodeURL(tc.input)
		if tc.hasError && err == nil {
			t.Errorf("Expected error for input %s, got none", tc.input)
		}
		if !tc.hasError && err != nil {
			t.Errorf("Expected no error for input %s, got %v", tc.input, err)
		}
		if !tc.hasError && result != tc.expected {
			t.Errorf("Expected %s for input %s, got %s", tc.expected, tc.input, result)
		}
	}
}

func TestHandler_ParseSearchQueries(t *testing.T) {
	handler, _, _ := setupTestHandler(t)

	testCases := []struct {
		args            []string
		expectedTargets []string
		expectedFilters []string
	}{
		{
			[]string{"test", "query"},
			[]string{"test", "query"},
			[]string{},
		},
		{
			[]string{"--familiarity=3", "test"},
			[]string{"test"},
			[]string{"--familiarity=3"},
		},
		{
			[]string{"--familiarity=3", "--importance=2", "test", "query"},
			[]string{"test", "query"},
			[]string{"--familiarity=3", "--importance=2"},
		},
		{
			[]string{"--familiarity=3"},
			[]string{},
			[]string{"--familiarity=3"},
		},
		{
			[]string{},
			[]string{},
			[]string{},
		},
	}

	for _, tc := range testCases {
		targets, filters := handler.parseSearchQueries(tc.args)

		if len(targets) != len(tc.expectedTargets) {
			t.Errorf("Expected %d targets, got %d", len(tc.expectedTargets), len(targets))
		}

		if len(filters) != len(tc.expectedFilters) {
			t.Errorf("Expected %d filters, got %d", len(tc.expectedFilters), len(filters))
		}

		for i, target := range targets {
			if target != tc.expectedTargets[i] {
				t.Errorf("Expected target %s, got %s", tc.expectedTargets[i], target)
			}
		}

		for i, filter := range filters {
			if filter != tc.expectedFilters[i] {
				t.Errorf("Expected filter %s, got %s", tc.expectedFilters[i], filter)
			}
		}
	}
}

func TestHandler_ParseFilterArgs(t *testing.T) {
	handler, _, _ := setupTestHandler(t)

	testCases := []struct {
		args     []string
		hasError bool
	}{
		{[]string{"--familiarity=3"}, false},
		{[]string{"--importance=2"}, false},
		{[]string{"--review-count=5"}, false},
		{[]string{"--due-only"}, false},
		{[]string{"--familiarity=invalid"}, true},
		{[]string{"--importance=invalid"}, true},
		{[]string{"--review-count=invalid"}, true},
		{[]string{"--unknown=value"}, false}, // Should be ignored
	}

	for _, tc := range testCases {
		_, err := handler.parseFilterArgs(tc.args)
		if tc.hasError && err == nil {
			t.Errorf("Expected error for args %v, got none", tc.args)
		}
		if !tc.hasError && err != nil {
			t.Errorf("Expected no error for args %v, got %v", tc.args, err)
		}
	}
}

func TestHandler_HandleUnknown(t *testing.T) {
	handler, mockIO, _ := setupTestHandler(t)

	handler.HandleUnknown("unknown_command")

	// Verify that warning was printed
	found := false
	for _, call := range mockIO.writeCalls {
		if call == "PrintfColored" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected PrintfColored to be called for unknown command")
	}
}

func TestHandler_HandleHelp(t *testing.T) {
	handler, mockIO, _ := setupTestHandler(t)

	handler.HandleHelp()

	// Verify that help was printed
	found := false
	for _, call := range mockIO.writeCalls {
		if call == "PrintlnColored" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected PrintlnColored to be called for help")
	}
}

func TestHandler_HandleClear(t *testing.T) {
	handler, mockIO, _ := setupTestHandler(t)

	handler.HandleClear()

	// Verify that clear was called
	found := false
	for _, call := range mockIO.writeCalls {
		if call == "Println" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected Println to be called for clear")
	}
}

func TestHandler_HandleQuit(t *testing.T) {
	handler, mockIO, _ := setupTestHandler(t)

	handler.HandleQuit()

	// Verify that quit message was printed
	output := mockIO.output.String()
	if !strings.Contains(output, "Goodbye!") {
		t.Error("Expected 'Goodbye!' message to be printed")
	}
}

func TestHandler_HandleHistory_Success(t *testing.T) {
	handler, mockIO, mockUseCase := setupTestHandler(t)

	// Set up history data
	mockUseCase.shouldError = false

	handler.HandleHistory()

	// Verify that history was displayed
	output := mockIO.output.String()
	if !strings.Contains(output, "Action History") {
		t.Error("Expected history header to be displayed")
	}
	if !strings.Contains(output, "test-question") {
		t.Error("Expected history entry to be displayed")
	}
}

func TestHandler_HandleHistory_Error(t *testing.T) {
	handler, mockIO, mockUseCase := setupTestHandler(t)

	// Set up error
	mockUseCase.shouldError = true
	mockUseCase.errorToReturn = errors.New("test error")

	handler.HandleHistory()

	// Verify that error was printed
	found := false
	for _, call := range mockIO.writeCalls {
		if call == "PrintError" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected PrintError to be called for history error")
	}
}

func TestHandler_GetQuestionsPage(t *testing.T) {
	handler, _, _ := setupTestHandler(t)

	// Create test questions
	questions := []core.Question{
		{ID: 1, URL: "https://leetcode.com/problems/test1"},
		{ID: 2, URL: "https://leetcode.com/problems/test2"},
		{ID: 3, URL: "https://leetcode.com/problems/test3"},
		{ID: 4, URL: "https://leetcode.com/problems/test4"},
		{ID: 5, URL: "https://leetcode.com/problems/test5"},
	}

	// Test first page
	results, totalPages, err := handler.getQuestionsPage(questions, 0)
	if err != nil {
		t.Fatalf("Failed to get first page: %v", err)
	}

	if len(results) != 3 { // Page size is 3 from test config
		t.Errorf("Expected 3 questions on first page, got %d", len(results))
	}

	if totalPages != 2 { // 5 questions with page size 3 = 2 pages
		t.Errorf("Expected 2 total pages, got %d", totalPages)
	}

	// Test with more questions to test pagination
	moreQuestions := []core.Question{
		{ID: 1, URL: "https://leetcode.com/problems/test1"},
		{ID: 2, URL: "https://leetcode.com/problems/test2"},
		{ID: 3, URL: "https://leetcode.com/problems/test3"},
		{ID: 4, URL: "https://leetcode.com/problems/test4"},
		{ID: 5, URL: "https://leetcode.com/problems/test5"},
		{ID: 6, URL: "https://leetcode.com/problems/test6"},
		{ID: 7, URL: "https://leetcode.com/problems/test7"},
	}

	// Test first page with more questions
	results, totalPages, err = handler.getQuestionsPage(moreQuestions, 0)
	if err != nil {
		t.Fatalf("Failed to get first page with more questions: %v", err)
	}

	if len(results) != 3 { // Page size is 3
		t.Errorf("Expected 3 questions on first page, got %d", len(results))
	}

	if totalPages != 3 { // 7 questions with page size 3 = 3 pages
		t.Errorf("Expected 3 total pages, got %d", totalPages)
	}

	// Test second page
	results, totalPages, err = handler.getQuestionsPage(moreQuestions, 1)
	if err != nil {
		t.Fatalf("Failed to get second page: %v", err)
	}

	if len(results) != 3 { // Page size is 3
		t.Errorf("Expected 3 questions on second page, got %d", len(results))
	}

	// Test invalid page number
	_, _, err = handler.getQuestionsPage(questions, -1)
	if err == nil {
		t.Error("Expected error for invalid page number")
	}

	// Test page number too high (page 2 for 5 questions with page size 3 = only pages 0,1 exist)
	_, _, err = handler.getQuestionsPage(questions, 2)
	if err == nil {
		t.Error("Expected error for page number too high")
	}

	// Test empty questions
	emptyQuestions := []core.Question{}
	results, totalPages, err = handler.getQuestionsPage(emptyQuestions, 0)
	if err != nil {
		t.Fatalf("Failed to get page for empty questions: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results for empty questions, got %d", len(results))
	}

	if totalPages != 0 {
		t.Errorf("Expected 0 total pages for empty questions, got %d", totalPages)
	}
}
