// Package urlparser provides URL parsing and normalization for supported DSA platforms.
package urlparser

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/eannchen/leetsolv/core"
	"github.com/eannchen/leetsolv/internal/errs"
)

// platformParser defines how to parse URLs for a specific platform
type platformParser struct {
	platform     core.Platform
	pathPrefix   string
	pathRegex    *regexp.Regexp
	normalizeURL func(slug string) string
}

// newLeetCodeParser creates a parser for LeetCode URLs
func newLeetCodeParser(host string) platformParser {
	return platformParser{
		platform:   core.PlatformLeetCode,
		pathPrefix: "/problems/",
		pathRegex:  regexp.MustCompile(`^/problems/([^/]+)`),
		normalizeURL: func(slug string) string {
			return "https://" + host + "/problems/" + slug + "/"
		},
	}
}

// newHackerRankParser creates a parser for HackerRank URLs
func newHackerRankParser(host string) platformParser {
	return platformParser{
		platform:   core.PlatformHackerRank,
		pathPrefix: "/challenges/",
		pathRegex:  regexp.MustCompile(`^/challenges/([^/]+)`),
		normalizeURL: func(slug string) string {
			return "https://" + host + "/challenges/" + slug + "/problem"
		},
	}
}

// platformParsers maps hostnames to their parser configuration
var platformParsers = map[string]platformParser{
	"leetcode.com":       newLeetCodeParser("leetcode.com"),
	"hackerrank.com":     newHackerRankParser("hackerrank.com"),
	"www.hackerrank.com": newHackerRankParser("www.hackerrank.com"),
}

// Parse normalizes and validates a URL from any supported platform.
// Returns a ParsedURL with the platform, normalized URL, and problem slug.
//
// Example inputs and outputs:
//   - "https://leetcode.com/problems/two-sum" → {LeetCode, "https://leetcode.com/problems/two-sum/", "two-sum"}
//   - "https://hackerrank.com/challenges/solve-me-first/problem" → {HackerRank, "https://hackerrank.com/challenges/solve-me-first/problem", "solve-me-first"}
func Parse(inputURL string) (*core.ParsedURL, error) {
	// Parse and validate URL structure
	parsedURL, err := url.Parse(strings.TrimSpace(inputURL))
	if err != nil {
		return nil, errs.ErrInvalidURLFormat
	}

	// O(1) lookup by host
	parser, found := platformParsers[parsedURL.Host]
	if !found || !strings.HasPrefix(parsedURL.Path, parser.pathPrefix) {
		return nil, errs.ErrUnsupportedPlatform
	}

	// Extract problem slug using platform-specific regex
	matches := parser.pathRegex.FindStringSubmatch(parsedURL.Path)
	if len(matches) != 2 {
		return nil, errs.ErrInvalidProblemURLFormat
	}

	slug := strings.TrimSpace(matches[1])
	if slug == "" {
		return nil, errs.ErrInvalidProblemURLFormat
	}

	// Build and return the parsed result
	return &core.ParsedURL{
		Platform:      parser.platform,
		NormalizedURL: parser.normalizeURL(slug),
		ProblemSlug:   slug,
	}, nil
}
