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
	host         string
	pathPrefix   string
	pathRegex    *regexp.Regexp
	normalizeURL func(slug string) string
}

var platformParsers = []platformParser{
	{
		host:       "leetcode.com",
		pathPrefix: "/problems/",
		pathRegex:  regexp.MustCompile(`^/problems/([^/]+)`),
		normalizeURL: func(slug string) string {
			return "https://leetcode.com/problems/" + slug + "/"
		},
	},
	{
		host:       "hackerrank.com",
		pathPrefix: "/challenges/",
		pathRegex:  regexp.MustCompile(`^/challenges/([^/]+)`),
		normalizeURL: func(slug string) string {
			return "https://hackerrank.com/challenges/" + slug + "/problem"
		},
	},
	{
		host:       "www.hackerrank.com",
		pathPrefix: "/challenges/",
		pathRegex:  regexp.MustCompile(`^/challenges/([^/]+)`),
		normalizeURL: func(slug string) string {
			return "https://www.hackerrank.com/challenges/" + slug + "/problem"
		},
	},
}

// hostToPlatform maps hostnames to Platform type
var hostToPlatform = map[string]core.Platform{
	"leetcode.com":       core.PlatformLeetCode,
	"hackerrank.com":     core.PlatformHackerRank,
	"www.hackerrank.com": core.PlatformHackerRank,
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

	// Find matching platform by host and path prefix
	for _, parser := range platformParsers {
		if parsedURL.Host == parser.host && strings.HasPrefix(parsedURL.Path, parser.pathPrefix) {
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
			platform := hostToPlatform[parsedURL.Host]
			return &core.ParsedURL{
				Platform:      platform,
				NormalizedURL: parser.normalizeURL(slug),
				ProblemSlug:   slug,
			}, nil
		}
	}

	// No matching platform found
	return nil, errs.ErrUnsupportedPlatform
}

// SupportedPlatforms returns a list of supported platform names
func SupportedPlatforms() []string {
	return []string{"LeetCode", "HackerRank"}
}
