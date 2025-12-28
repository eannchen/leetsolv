package urlparser

import (
	"testing"

	"github.com/eannchen/leetsolv/core"
)

func TestParse_LeetCode(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expectedURL   string
		expectedSlug  string
		expectedError bool
	}{
		{
			name:          "basic LeetCode URL",
			input:         "https://leetcode.com/problems/two-sum/",
			expectedURL:   "https://leetcode.com/problems/two-sum/",
			expectedSlug:  "two-sum",
			expectedError: false,
		},
		{
			name:          "LeetCode URL without trailing slash",
			input:         "https://leetcode.com/problems/two-sum",
			expectedURL:   "https://leetcode.com/problems/two-sum/",
			expectedSlug:  "two-sum",
			expectedError: false,
		},
		{
			name:          "LeetCode URL with solution path",
			input:         "https://leetcode.com/problems/two-sum/solution/",
			expectedURL:   "https://leetcode.com/problems/two-sum/",
			expectedSlug:  "two-sum",
			expectedError: false,
		},
		{
			name:          "LeetCode URL with discuss path",
			input:         "https://leetcode.com/problems/two-sum/discuss/",
			expectedURL:   "https://leetcode.com/problems/two-sum/",
			expectedSlug:  "two-sum",
			expectedError: false,
		},
		{
			name:          "LeetCode URL empty problem",
			input:         "https://leetcode.com/problems/",
			expectedError: true,
		},
		{
			name:          "LeetCode URL no problem path",
			input:         "https://leetcode.com/",
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Parse(tc.input)
			if tc.expectedError {
				if err == nil {
					t.Errorf("expected error for input %s, got none", tc.input)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error for input %s: %v", tc.input, err)
				return
			}
			if result.Platform != core.PlatformLeetCode {
				t.Errorf("expected platform LeetCode, got %s", result.Platform)
			}
			if result.NormalizedURL != tc.expectedURL {
				t.Errorf("expected URL %s, got %s", tc.expectedURL, result.NormalizedURL)
			}
			if result.ProblemSlug != tc.expectedSlug {
				t.Errorf("expected slug %s, got %s", tc.expectedSlug, result.ProblemSlug)
			}
		})
	}
}

func TestParse_HackerRank(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expectedURL   string
		expectedSlug  string
		expectedError bool
	}{
		{
			name:          "basic HackerRank URL",
			input:         "https://hackerrank.com/challenges/solve-me-first/problem",
			expectedURL:   "https://hackerrank.com/challenges/solve-me-first/problem",
			expectedSlug:  "solve-me-first",
			expectedError: false,
		},
		{
			name:          "HackerRank URL with www",
			input:         "https://www.hackerrank.com/challenges/solve-me-first/problem",
			expectedURL:   "https://www.hackerrank.com/challenges/solve-me-first/problem",
			expectedSlug:  "solve-me-first",
			expectedError: false,
		},
		{
			name:          "HackerRank URL without problem suffix",
			input:         "https://hackerrank.com/challenges/simple-array-sum",
			expectedURL:   "https://hackerrank.com/challenges/simple-array-sum/problem",
			expectedSlug:  "simple-array-sum",
			expectedError: false,
		},
		{
			name:          "HackerRank URL with domain path",
			input:         "https://www.hackerrank.com/challenges/breaking-the-records/problem?isFullScreen=true",
			expectedURL:   "https://www.hackerrank.com/challenges/breaking-the-records/problem",
			expectedSlug:  "breaking-the-records",
			expectedError: false,
		},
		{
			name:          "HackerRank URL empty challenge",
			input:         "https://hackerrank.com/challenges/",
			expectedError: true,
		},
		{
			name:          "HackerRank URL no challenges path",
			input:         "https://hackerrank.com/",
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Parse(tc.input)
			if tc.expectedError {
				if err == nil {
					t.Errorf("expected error for input %s, got none", tc.input)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error for input %s: %v", tc.input, err)
				return
			}
			if result.Platform != core.PlatformHackerRank {
				t.Errorf("expected platform HackerRank, got %s", result.Platform)
			}
			if result.NormalizedURL != tc.expectedURL {
				t.Errorf("expected URL %s, got %s", tc.expectedURL, result.NormalizedURL)
			}
			if result.ProblemSlug != tc.expectedSlug {
				t.Errorf("expected slug %s, got %s", tc.expectedSlug, result.ProblemSlug)
			}
		})
	}
}

func TestParse_UnsupportedPlatform(t *testing.T) {
	testCases := []string{
		"https://google.com/problems/test",
		"https://codeforces.com/problemset/problem/1/A",
		"https://example.com/",
		"invalid-url",
		"",
	}

	for _, input := range testCases {
		t.Run(input, func(t *testing.T) {
			_, err := Parse(input)
			if err == nil {
				t.Errorf("expected error for unsupported URL %s, got none", input)
			}
		})
	}
}

func TestSupportedPlatforms(t *testing.T) {
	platforms := SupportedPlatforms()
	if len(platforms) != 2 {
		t.Errorf("expected 2 supported platforms, got %d", len(platforms))
	}

	expected := map[string]bool{"LeetCode": true, "HackerRank": true}
	for _, p := range platforms {
		if !expected[p] {
			t.Errorf("unexpected platform: %s", p)
		}
	}
}
