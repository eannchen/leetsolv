package urlparser

import (
	"errors"
	"testing"

	"github.com/eannchen/leetsolv/core"
	"github.com/eannchen/leetsolv/internal/errs"
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
			name:          "LeetCode URL with query parameters",
			input:         "https://leetcode.com/problems/two-sum?envType=study-plan",
			expectedURL:   "https://leetcode.com/problems/two-sum/",
			expectedSlug:  "two-sum",
			expectedError: false,
		},
		{
			name:          "LeetCode URL with fragment",
			input:         "https://leetcode.com/problems/two-sum#description",
			expectedURL:   "https://leetcode.com/problems/two-sum/",
			expectedSlug:  "two-sum",
			expectedError: false,
		},
		{
			name:          "LeetCode URL with leading whitespace",
			input:         "  https://leetcode.com/problems/two-sum/",
			expectedURL:   "https://leetcode.com/problems/two-sum/",
			expectedSlug:  "two-sum",
			expectedError: false,
		},
		{
			name:          "LeetCode URL with trailing whitespace",
			input:         "https://leetcode.com/problems/two-sum/  ",
			expectedURL:   "https://leetcode.com/problems/two-sum/",
			expectedSlug:  "two-sum",
			expectedError: false,
		},
		{
			name:          "LeetCode URL with numeric slug",
			input:         "https://leetcode.com/problems/3sum/",
			expectedURL:   "https://leetcode.com/problems/3sum/",
			expectedSlug:  "3sum",
			expectedError: false,
		},
		{
			name:          "LeetCode URL with HTTP scheme",
			input:         "http://leetcode.com/problems/two-sum/",
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
			name:          "HackerRank URL with query params",
			input:         "https://www.hackerrank.com/challenges/breaking-the-records/problem?isFullScreen=true",
			expectedURL:   "https://www.hackerrank.com/challenges/breaking-the-records/problem",
			expectedSlug:  "breaking-the-records",
			expectedError: false,
		},
		{
			name:          "HackerRank URL with fragment",
			input:         "https://hackerrank.com/challenges/solve-me-first#problem",
			expectedURL:   "https://hackerrank.com/challenges/solve-me-first/problem",
			expectedSlug:  "solve-me-first",
			expectedError: false,
		},
		{
			name:          "HackerRank URL with whitespace",
			input:         "  https://hackerrank.com/challenges/solve-me-first/problem  ",
			expectedURL:   "https://hackerrank.com/challenges/solve-me-first/problem",
			expectedSlug:  "solve-me-first",
			expectedError: false,
		},
		{
			name:          "HackerRank URL with HTTP scheme",
			input:         "http://hackerrank.com/challenges/solve-me-first/problem",
			expectedURL:   "https://hackerrank.com/challenges/solve-me-first/problem",
			expectedSlug:  "solve-me-first",
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

func TestParse_ErrorTypes(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expectedErr error
	}{
		{
			name:        "unsupported platform returns ErrUnsupportedPlatform",
			input:       "https://codeforces.com/problem/123",
			expectedErr: errs.ErrUnsupportedPlatform,
		},
		{
			name:        "empty URL returns ErrUnsupportedPlatform",
			input:       "",
			expectedErr: errs.ErrUnsupportedPlatform,
		},
		{
			name:        "supported host but wrong path returns ErrUnsupportedPlatform",
			input:       "https://leetcode.com/contest/weekly-123",
			expectedErr: errs.ErrUnsupportedPlatform,
		},
		{
			name:        "empty problem slug returns ErrInvalidProblemURLFormat",
			input:       "https://leetcode.com/problems/",
			expectedErr: errs.ErrInvalidProblemURLFormat,
		},
		{
			name:        "empty challenge slug returns ErrInvalidProblemURLFormat",
			input:       "https://hackerrank.com/challenges/",
			expectedErr: errs.ErrInvalidProblemURLFormat,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.input)
			if err == nil {
				t.Errorf("expected error for input %s, got none", tc.input)
				return
			}
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

func TestParse_HostVariations(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expectedError bool
	}{
		{
			name:          "uppercase host should fail (URL parsing is case-sensitive for host lookup)",
			input:         "https://LEETCODE.COM/problems/two-sum",
			expectedError: true,
		},
		{
			name:          "mixed case host should fail",
			input:         "https://LeetCode.com/problems/two-sum",
			expectedError: true,
		},
		{
			name:          "www.leetcode.com is not supported",
			input:         "https://www.leetcode.com/problems/two-sum",
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.input)
			if tc.expectedError && err == nil {
				t.Errorf("expected error for input %s, got none", tc.input)
			}
			if !tc.expectedError && err != nil {
				t.Errorf("unexpected error for input %s: %v", tc.input, err)
			}
		})
	}
}
