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
	tests := []struct {
		name    string
		err     error
		userMsg string
	}{
		{
			name:    "ErrInvalidPageNumber",
			err:     ErrInvalidPageNumber,
			userMsg: "Invalid page number",
		},
		{
			name:    "ErrInvalidURLFormat",
			err:     ErrInvalidURLFormat,
			userMsg: "Please provide a valid URL",
		},
		{
			name:    "ErrInvalidEmptyInput",
			err:     ErrInvalidEmptyInput,
			userMsg: "Please provide a valid input",
		},
		{
			name:    "ErrUnsupportedPlatform",
			err:     ErrUnsupportedPlatform,
			userMsg: "Unsupported platform or problem URL format. Supported: LeetCode, HackerRank",
		},
		{
			name:    "ErrInvalidProblemURLFormat",
			err:     ErrInvalidProblemURLFormat,
			userMsg: "Invalid problem URL format",
		},
		{
			name:    "ErrInvalidFamiliarityLevel",
			err:     ErrInvalidFamiliarityLevel,
			userMsg: "Please enter a familiarity level between 1 and 5",
		},
		{
			name:    "ErrInvalidImportanceLevel",
			err:     ErrInvalidImportanceLevel,
			userMsg: "Please enter an importance level between 1 and 4",
		},
		{
			name:    "ErrInvalidReviewCount",
			err:     ErrInvalidReviewCount,
			userMsg: "Please enter a valid review count",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Fatalf("%s should not be nil", tt.name)
			}

			codedErr, ok := tt.err.(*CodedError)
			if !ok {
				t.Fatalf("%s should be a CodedError", tt.name)
			}

			if codedErr.Kind != ValidationErrorKind {
				t.Errorf("%s kind is %s, expected %s", tt.name, codedErr.Kind, ValidationErrorKind)
			}

			if codedErr.UserMsg != tt.userMsg {
				t.Errorf("%s user message is %q, expected %q", tt.name, codedErr.UserMsg, tt.userMsg)
			}
		})
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
		ErrUnsupportedPlatform,
		ErrInvalidProblemURLFormat,
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
		ErrUnsupportedPlatform,
		ErrInvalidProblemURLFormat,
		ErrInvalidFamiliarityLevel,
		ErrInvalidImportanceLevel,
		ErrInvalidReviewCount,
	}

	for _, err := range errors {
		// This should not cause any compilation errors
		_ = err.Error()
	}
}
