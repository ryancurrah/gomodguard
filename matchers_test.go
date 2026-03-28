package gomodguard_test

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ryancurrah/gomodguard/v2"
)

func TestMatchers(t *testing.T) {
	tests := map[string]struct {
		matcher   gomodguard.Matcher
		input     string
		wantMatch bool
	}{
		"exact match": {
			matcher:   gomodguard.ExactMatcher{Target: "github.com/foo/bar"},
			input:     "github.com/foo/bar",
			wantMatch: true,
		},
		"exact no match different version": {
			matcher:   gomodguard.ExactMatcher{Target: "github.com/foo/bar"},
			input:     "github.com/foo/bar/v2",
			wantMatch: false,
		},
		"exact no match partial": {
			matcher:   gomodguard.ExactMatcher{Target: "github.com/foo/bar"},
			input:     "github.com/foo",
			wantMatch: false,
		},
		"prefix match subpath": {
			matcher:   gomodguard.PrefixMatcher{Prefix: "golang.org"},
			input:     "golang.org/x/mod",
			wantMatch: true,
		},
		"prefix match exact": {
			matcher:   gomodguard.PrefixMatcher{Prefix: "golang.org"},
			input:     "golang.org",
			wantMatch: true,
		},
		"prefix match case insensitive with whitespace": {
			matcher:   gomodguard.PrefixMatcher{Prefix: "golang.org"},
			input:     "  Golang.Org/x/tools  ",
			wantMatch: true,
		},
		"prefix no match different domain": {
			matcher:   gomodguard.PrefixMatcher{Prefix: "golang.org"},
			input:     "github.com/golang",
			wantMatch: false,
		},
		"regex match mod": {
			matcher:   gomodguard.RegexMatcher{Regex: regexp.MustCompile(`golang\.org/x/.*`)},
			input:     "golang.org/x/mod",
			wantMatch: true,
		},
		"regex match tools": {
			matcher:   gomodguard.RegexMatcher{Regex: regexp.MustCompile(`golang\.org/x/.*`)},
			input:     "golang.org/x/tools",
			wantMatch: true,
		},
		"regex no match": {
			matcher:   gomodguard.RegexMatcher{Regex: regexp.MustCompile(`golang\.org/x/.*`)},
			input:     "golang.org/dl",
			wantMatch: false,
		},
		"regex nil never matches": {
			matcher:   gomodguard.RegexMatcher{Regex: nil},
			input:     "anything",
			wantMatch: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.wantMatch, tc.matcher.Match(tc.input))
		})
	}
}

func TestMatchTypePrecedence(t *testing.T) {
	tests := map[string]struct {
		matchType      gomodguard.MatchType
		wantPrecedence int
	}{
		"exact":   {gomodguard.ExactMatch, 0},
		"prefix":  {gomodguard.PrefixMatch, 1},
		"regex":   {gomodguard.RegexMatch, 2},
		"default": {"", 0},
		"unknown": {"unknown", 0},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.wantPrecedence, tc.matchType.Precedence())
		})
	}

	// Verify ordering: exact < prefix < regex
	assert.Less(t, gomodguard.ExactMatch.Precedence(), gomodguard.PrefixMatch.Precedence())
	assert.Less(t, gomodguard.PrefixMatch.Precedence(), gomodguard.RegexMatch.Precedence())
}
