package handler

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/eannchen/leetsolv/core"
	"github.com/eannchen/leetsolv/internal/clock"
	"github.com/eannchen/leetsolv/internal/errs"
)

func TestIOHandler_FormatTimeAgo(t *testing.T) {
	// Fixed "now" time for deterministic tests
	fixedNow := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	mockClock := clock.NewMockClock(fixedNow)

	ioh := &IOHandlerImpl{
		Clock: mockClock,
	}

	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{
			name:     "just now (less than 1 minute)",
			input:    fixedNow.Add(-30 * time.Second),
			expected: "just now",
		},
		{
			name:     "1 minute ago",
			input:    fixedNow.Add(-1 * time.Minute),
			expected: "1 minute ago",
		},
		{
			name:     "multiple minutes ago",
			input:    fixedNow.Add(-30 * time.Minute),
			expected: "30 minutes ago",
		},
		{
			name:     "1 hour ago",
			input:    fixedNow.Add(-1 * time.Hour),
			expected: "1 hour ago",
		},
		{
			name:     "multiple hours ago",
			input:    fixedNow.Add(-5 * time.Hour),
			expected: "5 hours ago",
		},
		{
			name:     "1 day ago",
			input:    fixedNow.Add(-24 * time.Hour),
			expected: "1 day ago",
		},
		{
			name:     "multiple days ago",
			input:    fixedNow.Add(-72 * time.Hour),
			expected: "3 days ago",
		},
		{
			name:     "edge: exactly 59 seconds",
			input:    fixedNow.Add(-59 * time.Second),
			expected: "just now",
		},
		{
			name:     "edge: exactly 59 minutes",
			input:    fixedNow.Add(-59 * time.Minute),
			expected: "59 minutes ago",
		},
		{
			name:     "edge: exactly 23 hours",
			input:    fixedNow.Add(-23 * time.Hour),
			expected: "23 hours ago",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ioh.FormatTimeAgo(tt.input)
			if result != tt.expected {
				t.Errorf("FormatTimeAgo(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNewIOHandler(t *testing.T) {
	mockClock := clock.NewMockClock(time.Now())
	ioh := NewIOHandler(mockClock)

	if ioh == nil {
		t.Error("NewIOHandler should return non-nil")
	}
	if ioh.Clock != mockClock {
		t.Error("Clock should be set")
	}
	if ioh.Reader == nil {
		t.Error("Reader should be set")
	}
	if ioh.Writer == nil {
		t.Error("Writer should be set")
	}
}

func TestIOHandler_Println(t *testing.T) {
	var buf bytes.Buffer
	ioh := &IOHandlerImpl{Writer: &buf}

	ioh.Println("hello", "world")

	if !strings.Contains(buf.String(), "hello world") {
		t.Errorf("Expected 'hello world' in output, got %q", buf.String())
	}
}

func TestIOHandler_Printf(t *testing.T) {
	var buf bytes.Buffer
	ioh := &IOHandlerImpl{Writer: &buf}

	ioh.Printf("count: %d", 42)

	if buf.String() != "count: 42" {
		t.Errorf("Expected 'count: 42', got %q", buf.String())
	}
}

func TestIOHandler_PrintlnColored(t *testing.T) {
	var buf bytes.Buffer
	ioh := &IOHandlerImpl{Writer: &buf}

	ioh.PrintlnColored(ColorGreen, "success")

	output := buf.String()
	if !strings.Contains(output, "success") {
		t.Errorf("Expected 'success' in output, got %q", output)
	}
	if !strings.Contains(output, ColorGreen) {
		t.Errorf("Expected color code in output, got %q", output)
	}
}

func TestIOHandler_PrintfColored(t *testing.T) {
	var buf bytes.Buffer
	ioh := &IOHandlerImpl{Writer: &buf}

	ioh.PrintfColored(ColorRed, "error: %s", "failed")

	output := buf.String()
	if !strings.Contains(output, "error: failed") {
		t.Errorf("Expected 'error: failed' in output, got %q", output)
	}
	if !strings.Contains(output, ColorRed) {
		t.Errorf("Expected color code in output, got %q", output)
	}
}

func TestIOHandler_ReadLine(t *testing.T) {
	input := "test input\n"
	ioh := &IOHandlerImpl{
		Reader: strings.NewReader(input),
		Writer: &bytes.Buffer{},
	}

	scanner := bufio.NewScanner(ioh.Reader)
	result := ioh.ReadLine(scanner, "prompt> ")

	if result != "test input" {
		t.Errorf("Expected 'test input', got %q", result)
	}
}

func TestIOHandler_PrintQuestionBrief(t *testing.T) {
	var buf bytes.Buffer
	fixedNow := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	ioh := &IOHandlerImpl{
		Writer: &buf,
		Clock:  clock.NewMockClock(fixedNow),
	}

	q := &core.Question{
		ID:         1,
		URL:        "https://leetcode.com/problems/two-sum",
		Note:       "test note",
		NextReview: fixedNow.Add(24 * time.Hour),
	}

	ioh.PrintQuestionBrief(q)

	output := buf.String()
	if !strings.Contains(output, "[1]") {
		t.Errorf("Expected '[1]' in output, got %q", output)
	}
	if !strings.Contains(output, "two-sum") {
		t.Errorf("Expected 'two-sum' in output, got %q", output)
	}
	if !strings.Contains(output, "test note") {
		t.Errorf("Expected 'test note' in output, got %q", output)
	}
}

func TestIOHandler_PrintQuestionBrief_NoNote(t *testing.T) {
	var buf bytes.Buffer
	fixedNow := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	ioh := &IOHandlerImpl{
		Writer: &buf,
		Clock:  clock.NewMockClock(fixedNow),
	}

	q := &core.Question{
		ID:         1,
		URL:        "https://leetcode.com/problems/two-sum",
		Note:       "",
		NextReview: fixedNow.Add(24 * time.Hour),
	}

	ioh.PrintQuestionBrief(q)

	if !strings.Contains(buf.String(), "(none)") {
		t.Errorf("Expected '(none)' for empty note, got %q", buf.String())
	}
}

func TestIOHandler_PrintSuccess(t *testing.T) {
	var buf bytes.Buffer
	ioh := &IOHandlerImpl{Writer: &buf}

	ioh.PrintSuccess("operation complete")

	output := buf.String()
	if !strings.Contains(output, "operation complete") {
		t.Errorf("Expected 'operation complete' in output, got %q", output)
	}
	if !strings.Contains(output, "[✔]") {
		t.Errorf("Expected checkmark in output, got %q", output)
	}
}

func TestIOHandler_PrintCancel(t *testing.T) {
	var buf bytes.Buffer
	ioh := &IOHandlerImpl{Writer: &buf}

	ioh.PrintCancel("cancelled")

	output := buf.String()
	if !strings.Contains(output, "cancelled") {
		t.Errorf("Expected 'cancelled' in output, got %q", output)
	}
	if !strings.Contains(output, "[i]") {
		t.Errorf("Expected info icon in output, got %q", output)
	}
}

func TestIOHandler_PrintError(t *testing.T) {
	var buf bytes.Buffer
	ioh := &IOHandlerImpl{Writer: &buf}

	tests := []struct {
		name     string
		err      error
		contains string
	}{
		{
			name:     "nil error",
			err:      nil,
			contains: "", // Should not print anything
		},
		{
			name: "validation error",
			err: &errs.CodedError{
				Kind:    errs.ValidationErrorKind,
				UserMsg: "invalid input",
			},
			contains: "invalid input",
		},
		{
			name: "business error",
			err: &errs.CodedError{
				Kind:    errs.BusinessErrorKind,
				UserMsg: "not allowed",
			},
			contains: "not allowed",
		},
		{
			name: "system error",
			err: &errs.CodedError{
				Kind:         errs.SystemErrorKind,
				TechnicalMsg: "database failed",
			},
			contains: "database failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			ioh.PrintError(tt.err)
			if tt.contains != "" && !strings.Contains(buf.String(), tt.contains) {
				t.Errorf("Expected %q in output, got %q", tt.contains, buf.String())
			}
		})
	}
}

func TestIOHandler_PrintQuestionDetail(t *testing.T) {
	var buf bytes.Buffer
	fixedNow := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	ioh := &IOHandlerImpl{
		Writer: &buf,
		Clock:  clock.NewMockClock(fixedNow),
	}

	q := &core.Question{
		ID:           1,
		URL:          "https://leetcode.com/problems/two-sum",
		Note:         "test note",
		Familiarity:  core.Medium,
		Importance:   core.HighImportance,
		LastReviewed: fixedNow.Add(-24 * time.Hour),
		NextReview:   fixedNow.Add(24 * time.Hour),
		ReviewCount:  5,
		EaseFactor:   2.5,
		CreatedAt:    fixedNow.Add(-48 * time.Hour),
	}

	ioh.PrintQuestionDetail(q)

	output := buf.String()
	if !strings.Contains(output, "[1]") {
		t.Errorf("Expected '[1]' in output")
	}
	if !strings.Contains(output, "Familiarity") {
		t.Errorf("Expected 'Familiarity' in output")
	}
	if !strings.Contains(output, "Importance") {
		t.Errorf("Expected 'Importance' in output")
	}
}

func TestIOHandler_PrintQuestionDetail_Due(t *testing.T) {
	var buf bytes.Buffer
	fixedNow := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	ioh := &IOHandlerImpl{
		Writer: &buf,
		Clock:  clock.NewMockClock(fixedNow),
	}

	q := &core.Question{
		ID:           1,
		URL:          "https://leetcode.com/problems/two-sum",
		NextReview:   fixedNow.Add(-24 * time.Hour), // Due (in the past)
		LastReviewed: fixedNow.Add(-48 * time.Hour),
		CreatedAt:    fixedNow.Add(-72 * time.Hour),
	}

	ioh.PrintQuestionDetail(q)

	output := buf.String()
	if !strings.Contains(output, "(Due)") {
		t.Errorf("Expected '(Due)' in output for overdue question, got %q", output)
	}
}

func TestIOHandler_PrintQuestionUpsertDetail_NewQuestion(t *testing.T) {
	var buf bytes.Buffer
	fixedNow := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	ioh := &IOHandlerImpl{
		Writer: &buf,
		Clock:  clock.NewMockClock(fixedNow),
	}

	delta := &core.Delta{
		OldState: nil,
		NewState: &core.Question{
			ID:           1,
			URL:          "https://leetcode.com/problems/two-sum",
			Familiarity:  core.Medium,
			Importance:   core.MediumImportance,
			LastReviewed: fixedNow,
			NextReview:   fixedNow.Add(7 * 24 * time.Hour),
			ReviewCount:  1,
			EaseFactor:   2.0,
			CreatedAt:    fixedNow,
		},
	}

	ioh.PrintQuestionUpsertDetail(delta)

	output := buf.String()
	if !strings.Contains(output, "[1]") {
		t.Errorf("Expected '[1]' in output")
	}
}

func TestIOHandler_PrintQuestionUpsertDetail_UpdatedQuestion(t *testing.T) {
	var buf bytes.Buffer
	fixedNow := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	ioh := &IOHandlerImpl{
		Writer: &buf,
		Clock:  clock.NewMockClock(fixedNow),
	}

	delta := &core.Delta{
		OldState: &core.Question{
			ID:           1,
			Familiarity:  core.Easy,
			Importance:   core.LowImportance,
			LastReviewed: fixedNow.Add(-24 * time.Hour),
			NextReview:   fixedNow,
			ReviewCount:  1,
			EaseFactor:   2.0,
			CreatedAt:    fixedNow.Add(-48 * time.Hour),
		},
		NewState: &core.Question{
			ID:           1,
			URL:          "https://leetcode.com/problems/two-sum",
			Familiarity:  core.Medium,
			Importance:   core.HighImportance,
			LastReviewed: fixedNow,
			NextReview:   fixedNow.Add(7 * 24 * time.Hour),
			ReviewCount:  2,
			EaseFactor:   2.2,
			CreatedAt:    fixedNow.Add(-48 * time.Hour),
		},
	}

	ioh.PrintQuestionUpsertDetail(delta)

	output := buf.String()
	// Should show changes with arrows
	if !strings.Contains(output, "→") {
		t.Errorf("Expected '→' for changes in output")
	}
}
