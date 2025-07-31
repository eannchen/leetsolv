package handler

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"leetsolv/internal/errs"
)

func TestPrintError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: "",
		},
		{
			name:     "input error from usecase",
			err:      errs.Err400QuestionNotFound,
			expected: "⚠️ FAILED INPUT: question not found",
		},
		{
			name:     "system error from usecase",
			err:      errs.WrapInternalError(errors.New("database connection failed"), "failed to load data"),
			expected: "❌ SYSTEM ERROR: failed to load data",
		},
		{
			name:     "handler validation error - invalid URL",
			err:      errs.ErrInvalidURLFormat,
			expected: "⚠️ FAILED INPUT: Please provide a valid URL",
		},
		{
			name:     "handler validation error - invalid familiarity",
			err:      errs.ErrInvalidFamiliarityLevel,
			expected: "⚠️ FAILED INPUT: Please enter a familiarity level between 1 and 5",
		},
		{
			name:     "unknown error",
			err:      errors.New("unknown error type"),
			expected: "❌ Error: unknown error type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			ioh := &IOHandlerImpl{
				Writer: &buf,
			}

			ioh.PrintError(tt.err)

			output := buf.String()
			if tt.err == nil {
				if output != "" {
					t.Errorf("expected empty output for nil error, got: %q", output)
				}
				return
			}

			if !strings.Contains(output, tt.expected) {
				t.Errorf("expected output to contain %q, got: %q", tt.expected, output)
			}
		})
	}
}
