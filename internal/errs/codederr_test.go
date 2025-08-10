package errs

import (
	"errors"
	"testing"
)

func TestCodedError_Error(t *testing.T) {
	// Test with technical message
	originalErr := errors.New("database connection failed")
	codedErr := &CodedError{
		Err:          originalErr,
		Kind:         SystemErrorKind,
		TechnicalMsg: "Failed to connect to database",
	}

	expected := "SYSTEM ERROR: Failed to connect to database (database connection failed)"
	if codedErr.Error() != expected {
		t.Errorf("Error() returned %q, expected %q", codedErr.Error(), expected)
	}

	// Test without technical message
	codedErr.TechnicalMsg = ""
	expected = "SYSTEM ERROR: database connection failed"
	if codedErr.Error() != expected {
		t.Errorf("Error() returned %q, expected %q", codedErr.Error(), expected)
	}
}

func TestCodedError_UserMessage(t *testing.T) {
	// Test with user message
	originalErr := errors.New("validation failed")
	codedErr := &CodedError{
		Err:     originalErr,
		Kind:    ValidationErrorKind,
		UserMsg: "Please check your input and try again",
	}

	expected := "Please check your input and try again"
	if codedErr.UserMessage() != expected {
		t.Errorf("UserMessage() returned %q, expected %q", codedErr.UserMessage(), expected)
	}

	// Test without user message (fallback to technical)
	codedErr.UserMsg = ""
	expected = "validation failed"
	if codedErr.UserMessage() != expected {
		t.Errorf("UserMessage() returned %q, expected %q", codedErr.UserMessage(), expected)
	}
}

func TestCodedError_Unwrap(t *testing.T) {
	originalErr := errors.New("original error")
	codedErr := &CodedError{
		Err:  originalErr,
		Kind: BusinessErrorKind,
	}

	unwrapped := codedErr.Unwrap()
	if unwrapped != originalErr {
		t.Errorf("Unwrap() returned %v, expected %v", unwrapped, originalErr)
	}
}

func TestWrapInternalError(t *testing.T) {
	originalErr := errors.New("internal error")
	wrapped := WrapInternalError(originalErr, "System failure")

	codedErr, ok := wrapped.(*CodedError)
	if !ok {
		t.Fatal("WrapInternalError did not return a CodedError")
	}

	if codedErr.Err != originalErr {
		t.Errorf("Wrapped error is %v, expected %v", codedErr.Err, originalErr)
	}

	if codedErr.Kind != SystemErrorKind {
		t.Errorf("Error kind is %s, expected %s", codedErr.Kind, SystemErrorKind)
	}

	if codedErr.TechnicalMsg != "System failure" {
		t.Errorf("Technical message is %q, expected %q", codedErr.TechnicalMsg, "System failure")
	}

	if codedErr.UserMsg != "" {
		t.Errorf("User message should be empty, got %q", codedErr.UserMsg)
	}
}

func TestWrapValidationError(t *testing.T) {
	originalErr := errors.New("validation error")
	wrapped := WrapValidationError(originalErr, "Please fix the errors")

	codedErr, ok := wrapped.(*CodedError)
	if !ok {
		t.Fatal("WrapValidationError did not return a CodedError")
	}

	if codedErr.Err != originalErr {
		t.Errorf("Wrapped error is %v, expected %v", codedErr.Err, originalErr)
	}

	if codedErr.Kind != ValidationErrorKind {
		t.Errorf("Error kind is %s, expected %s", codedErr.Kind, ValidationErrorKind)
	}

	if codedErr.UserMsg != "Please fix the errors" {
		t.Errorf("User message is %q, expected %q", codedErr.UserMsg, "Please fix the errors")
	}

	if codedErr.TechnicalMsg != "" {
		t.Errorf("Technical message should be empty, got %q", codedErr.TechnicalMsg)
	}
}

func TestWrapBusinessError(t *testing.T) {
	originalErr := errors.New("business error")
	wrapped := WrapBusinessError(originalErr, "Operation cannot be completed")

	codedErr, ok := wrapped.(*CodedError)
	if !ok {
		t.Fatal("WrapBusinessError did not return a CodedError")
	}

	if codedErr.Err != originalErr {
		t.Errorf("Wrapped error is %v, expected %v", codedErr.Err, originalErr)
	}

	if codedErr.Kind != BusinessErrorKind {
		t.Errorf("Error kind is %s, expected %s", codedErr.Kind, BusinessErrorKind)
	}

	if codedErr.UserMsg != "Operation cannot be completed" {
		t.Errorf("User message is %q, expected %q", codedErr.UserMsg, "Operation cannot be completed")
	}

	if codedErr.TechnicalMsg != "" {
		t.Errorf("Technical message should be empty, got %q", codedErr.TechnicalMsg)
	}
}

func TestErrorKindConstants(t *testing.T) {
	// Test that error kind constants have expected values
	if ValidationErrorKind != "VALIDATION ERROR" {
		t.Errorf("ValidationErrorKind is %q, expected %q", ValidationErrorKind, "VALIDATION ERROR")
	}

	if BusinessErrorKind != "BUSINESS ERROR" {
		t.Errorf("BusinessErrorKind is %q, expected %q", BusinessErrorKind, "BUSINESS ERROR")
	}

	if SystemErrorKind != "SYSTEM ERROR" {
		t.Errorf("SystemErrorKind is %q, expected %q", SystemErrorKind, "SYSTEM ERROR")
	}
}

func TestCodedError_InterfaceCompliance(t *testing.T) {
	// Test that CodedError implements the error interface
	var _ error = &CodedError{}
}

func TestCodedError_WithNilError(t *testing.T) {
	// Test behavior with nil error
	codedErr := &CodedError{
		Err:     nil,
		Kind:    ValidationErrorKind,
		UserMsg: "Test message",
	}

	// Should not panic
	_ = codedErr.Error()
	_ = codedErr.UserMessage()
	_ = codedErr.Unwrap()
}
