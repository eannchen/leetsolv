package errs

import (
	"testing"
)

func TestBusinessLogicErrors(t *testing.T) {
	// Test ErrQuestionNotFound
	if ErrQuestionNotFound == nil {
		t.Error("ErrQuestionNotFound should not be nil")
	}

	codedErr, ok := ErrQuestionNotFound.(*CodedError)
	if !ok {
		t.Fatal("ErrQuestionNotFound should be a CodedError")
	}

	if codedErr.Kind != BusinessErrorKind {
		t.Errorf("ErrQuestionNotFound kind is %s, expected %s", codedErr.Kind, BusinessErrorKind)
	}

	if codedErr.UserMsg != "Question not found. Please check the ID or URL" {
		t.Errorf("ErrQuestionNotFound user message is %q, expected %q",
			codedErr.UserMsg, "Question not found. Please check the ID or URL")
	}

	// Test ErrNoQuestionsAvailable
	if ErrNoQuestionsAvailable == nil {
		t.Error("ErrNoQuestionsAvailable should not be nil")
	}

	codedErr, ok = ErrNoQuestionsAvailable.(*CodedError)
	if !ok {
		t.Fatal("ErrNoQuestionsAvailable should be a CodedError")
	}

	if codedErr.Kind != BusinessErrorKind {
		t.Errorf("ErrNoQuestionsAvailable kind is %s, expected %s", codedErr.Kind, BusinessErrorKind)
	}

	if codedErr.UserMsg != "No questions available yet" {
		t.Errorf("ErrNoQuestionsAvailable user message is %q, expected %q",
			codedErr.UserMsg, "No questions available yet")
	}

	// Test ErrNoActionsToUndo
	if ErrNoActionsToUndo == nil {
		t.Error("ErrNoActionsToUndo should not be nil")
	}

	codedErr, ok = ErrNoActionsToUndo.(*CodedError)
	if !ok {
		t.Fatal("ErrNoActionsToUndo should be a CodedError")
	}

	if codedErr.Kind != BusinessErrorKind {
		t.Errorf("ErrNoActionsToUndo kind is %s, expected %s", codedErr.Kind, BusinessErrorKind)
	}

	if codedErr.UserMsg != "No actions to undo" {
		t.Errorf("ErrNoActionsToUndo user message is %q, expected %q",
			codedErr.UserMsg, "No actions to undo")
	}
}

func TestValidationErrors(t *testing.T) {
	// Test ErrInvalidPageNumber
	if ErrInvalidPageNumber == nil {
		t.Error("ErrInvalidPageNumber should not be nil")
	}

	codedErr, ok := ErrInvalidPageNumber.(*CodedError)
	if !ok {
		t.Fatal("ErrInvalidPageNumber should be a CodedError")
	}

	if codedErr.Kind != ValidationErrorKind {
		t.Errorf("ErrInvalidPageNumber kind is %s, expected %s", codedErr.Kind, ValidationErrorKind)
	}

	if codedErr.UserMsg != "Invalid page number" {
		t.Errorf("ErrInvalidPageNumber user message is %q, expected %q",
			codedErr.UserMsg, "Invalid page number")
	}

	// Test ErrInvalidURLFormat
	if ErrInvalidURLFormat == nil {
		t.Error("ErrInvalidURLFormat should not be nil")
	}

	codedErr, ok = ErrInvalidURLFormat.(*CodedError)
	if !ok {
		t.Fatal("ErrInvalidURLFormat should be a CodedError")
	}

	if codedErr.Kind != ValidationErrorKind {
		t.Errorf("ErrInvalidURLFormat kind is %s, expected %s", codedErr.Kind, ValidationErrorKind)
	}

	if codedErr.UserMsg != "Please provide a valid URL" {
		t.Errorf("ErrInvalidURLFormat user message is %q, expected %q",
			codedErr.UserMsg, "Please provide a valid URL")
	}

	// Test ErrInvalidEmptyInput
	if ErrInvalidEmptyInput == nil {
		t.Error("ErrInvalidEmptyInput should not be nil")
	}

	codedErr, ok = ErrInvalidEmptyInput.(*CodedError)
	if !ok {
		t.Fatal("ErrInvalidEmptyInput should be a CodedError")
	}

	if codedErr.Kind != ValidationErrorKind {
		t.Errorf("ErrInvalidEmptyInput kind is %s, expected %s", codedErr.Kind, ValidationErrorKind)
	}

	if codedErr.UserMsg != "Please provide a valid input" {
		t.Errorf("ErrInvalidEmptyInput user message is %q, expected %q",
			codedErr.UserMsg, "Please provide a valid input")
	}

	// Test ErrInvalidLeetCodeURL
	if ErrInvalidLeetCodeURL == nil {
		t.Error("ErrInvalidLeetCodeURL should not be nil")
	}

	codedErr, ok = ErrInvalidLeetCodeURL.(*CodedError)
	if !ok {
		t.Fatal("ErrInvalidLeetCodeURL should be a CodedError")
	}

	if codedErr.Kind != ValidationErrorKind {
		t.Errorf("ErrInvalidLeetCodeURL kind is %s, expected %s", codedErr.Kind, ValidationErrorKind)
	}

	if codedErr.UserMsg != "Please provide a valid LeetCode problem URL" {
		t.Errorf("ErrInvalidLeetCodeURL user message is %q, expected %q",
			codedErr.UserMsg, "Please provide a valid LeetCode problem URL")
	}

	// Test ErrInvalidLeetCodeURLFormat
	if ErrInvalidLeetCodeURLFormat == nil {
		t.Error("ErrInvalidLeetCodeURLFormat should not be nil")
	}

	codedErr, ok = ErrInvalidLeetCodeURLFormat.(*CodedError)
	if !ok {
		t.Fatal("ErrInvalidLeetCodeURLFormat should be a CodedError")
	}

	if codedErr.Kind != ValidationErrorKind {
		t.Errorf("ErrInvalidLeetCodeURLFormat kind is %s, expected %s", codedErr.Kind, ValidationErrorKind)
	}

	if codedErr.UserMsg != "Please provide a valid LeetCode problem URL format" {
		t.Errorf("ErrInvalidLeetCodeURLFormat user message is %q, expected %q",
			codedErr.UserMsg, "Please provide a valid LeetCode problem URL format")
	}

	// Test ErrInvalidFamiliarityLevel
	if ErrInvalidFamiliarityLevel == nil {
		t.Error("ErrInvalidFamiliarityLevel should not be nil")
	}

	codedErr, ok = ErrInvalidFamiliarityLevel.(*CodedError)
	if !ok {
		t.Fatal("ErrInvalidFamiliarityLevel should be a CodedError")
	}

	if codedErr.Kind != ValidationErrorKind {
		t.Errorf("ErrInvalidFamiliarityLevel kind is %s, expected %s", codedErr.Kind, ValidationErrorKind)
	}

	if codedErr.UserMsg != "Please enter a familiarity level between 1 and 5" {
		t.Errorf("ErrInvalidFamiliarityLevel user message is %q, expected %q",
			codedErr.UserMsg, "Please enter a familiarity level between 1 and 5")
	}

	// Test ErrInvalidImportanceLevel
	if ErrInvalidImportanceLevel == nil {
		t.Error("ErrInvalidImportanceLevel should not be nil")
	}

	codedErr, ok = ErrInvalidImportanceLevel.(*CodedError)
	if !ok {
		t.Fatal("ErrInvalidImportanceLevel should be a CodedError")
	}

	if codedErr.Kind != ValidationErrorKind {
		t.Errorf("ErrInvalidImportanceLevel kind is %s, expected %s", codedErr.Kind, ValidationErrorKind)
	}

	if codedErr.UserMsg != "Please enter an importance level between 1 and 4" {
		t.Errorf("ErrInvalidImportanceLevel user message is %q, expected %q",
			codedErr.UserMsg, "Please enter an importance level between 1 and 4")
	}

	// Test ErrInvalidReviewCount
	if ErrInvalidReviewCount == nil {
		t.Error("ErrInvalidReviewCount should not be nil")
	}

	codedErr, ok = ErrInvalidReviewCount.(*CodedError)
	if !ok {
		t.Fatal("ErrInvalidReviewCount should be a CodedError")
	}

	if codedErr.Kind != ValidationErrorKind {
		t.Errorf("ErrInvalidReviewCount kind is %s, expected %s", codedErr.Kind, ValidationErrorKind)
	}

	if codedErr.UserMsg != "Please enter a valid review count" {
		t.Errorf("ErrInvalidReviewCount user message is %q, expected %q",
			codedErr.UserMsg, "Please enter a valid review count")
	}
}

func TestErrorUnwrapping(t *testing.T) {
	// Test that all errors can be unwrapped to get the original error
	errors := []error{
		ErrQuestionNotFound,
		ErrNoQuestionsAvailable,
		ErrNoActionsToUndo,
		ErrInvalidPageNumber,
		ErrInvalidURLFormat,
		ErrInvalidEmptyInput,
		ErrInvalidLeetCodeURL,
		ErrInvalidLeetCodeURLFormat,
		ErrInvalidFamiliarityLevel,
		ErrInvalidImportanceLevel,
		ErrInvalidReviewCount,
	}

	for _, err := range errors {
		codedErr, ok := err.(*CodedError)
		if !ok {
			t.Errorf("Error %v is not a CodedError", err)
			continue
		}

		unwrapped := codedErr.Unwrap()
		if unwrapped == nil {
			t.Errorf("Unwrapped error for %v is nil", err)
		}
	}
}

func TestErrorInterfaceCompliance(t *testing.T) {
	// Test that all errors implement the error interface
	errors := []error{
		ErrQuestionNotFound,
		ErrNoQuestionsAvailable,
		ErrNoActionsToUndo,
		ErrInvalidPageNumber,
		ErrInvalidURLFormat,
		ErrInvalidEmptyInput,
		ErrInvalidLeetCodeURL,
		ErrInvalidLeetCodeURLFormat,
		ErrInvalidFamiliarityLevel,
		ErrInvalidImportanceLevel,
		ErrInvalidReviewCount,
	}

	for _, err := range errors {
		// This should not cause any compilation errors
		_ = err.Error()
	}
}
