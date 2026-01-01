package handler

import (
	"testing"
	"time"

	"github.com/eannchen/leetsolv/internal/clock"
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
