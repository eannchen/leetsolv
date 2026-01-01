package core

import "testing"

func TestPlatformString(t *testing.T) {
	tests := []struct {
		platform Platform
		expected string
	}{
		{PlatformLeetCode, "LeetCode"},
		{PlatformHackerRank, "HackerRank"},
		{Platform("unknown"), "unknown"},
	}

	for _, tt := range tests {
		t.Run(string(tt.platform), func(t *testing.T) {
			if got := tt.platform.String(); got != tt.expected {
				t.Errorf("Platform.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestActionTypeString(t *testing.T) {
	tests := []struct {
		action   ActionType
		expected string
	}{
		{ActionAdd, "Add"},
		{ActionUpdate, "Update"},
		{ActionDelete, "Delete"},
		{ActionType("unknown"), ""},
	}

	for _, tt := range tests {
		t.Run(string(tt.action), func(t *testing.T) {
			if got := tt.action.String(); got != tt.expected {
				t.Errorf("ActionType.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestActionTypePastTenseString(t *testing.T) {
	tests := []struct {
		action   ActionType
		expected string
	}{
		{ActionAdd, "Added"},
		{ActionUpdate, "Updated"},
		{ActionDelete, "Deleted"},
		{ActionType("unknown"), ""},
	}

	for _, tt := range tests {
		t.Run(string(tt.action), func(t *testing.T) {
			if got := tt.action.PastTenseString(); got != tt.expected {
				t.Errorf("ActionType.PastTenseString() = %q, want %q", got, tt.expected)
			}
		})
	}
}
