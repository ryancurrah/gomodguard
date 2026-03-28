package gomodguard_test

import (
	"os"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ryancurrah/gomodguard/v2"
)

func mustConstraint(t *testing.T, c string) *semver.Constraints {
	t.Helper()

	cs, err := semver.NewConstraint(c)
	require.NoError(t, err)

	return cs
}

func TestProcessorNewProcessor(t *testing.T) {
	_, err := gomodguard.NewProcessor(&gomodguard.Configuration{
		Allowed: gomodguard.Allowed{
			"github.com/foo/bar": gomodguard.AllowedRule{},
		},
	})
	require.NoError(t, err)
}

func TestProcessorNewProcessorUnknownMatchType(t *testing.T) {
	_, err := gomodguard.NewProcessor(&gomodguard.Configuration{
		Blocked: gomodguard.Blocked{
			"github.com/foo/bar": gomodguard.BlockedRule{
				MatchType: "prefx",
			},
		},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown match_type")
}

func processFiles(t *testing.T, config *gomodguard.Configuration) []string {
	t.Helper()

	wd, err := os.Getwd()
	require.NoError(t, err)

	processor, err := gomodguard.NewProcessor(config)
	require.NoError(t, err)

	processor.SetBlockedModules()

	filteredFiles := gomodguard.Find(wd, false, []string{"./..."})
	results := processor.ProcessFiles(filteredFiles)

	reasons := make([]string, 0, len(results))
	for _, r := range results {
		reasons = append(reasons, r.String())
	}

	return reasons
}

func TestProcessorProcessFiles(t *testing.T) { //nolint:funlen
	tests := map[string]struct {
		exampleDir     string
		config         *gomodguard.Configuration
		wantReasons    []string
		notWantReasons []string
		wantEmpty      bool
	}{
		"current module is a recommendation - not blocked": {
			exampleDir: "examples/alloptions",
			config: &gomodguard.Configuration{
				Blocked: gomodguard.Blocked{
					"github.com/gofrs/uuid": gomodguard.BlockedRule{
						Recommendations: []string{"github.com/ryancurrah/gomodguard/examples/alloptions"},
						Reason:          "should be skipped because current module is a recommendation.",
					},
				},
			},
			notWantReasons: []string{"github.com/gofrs/uuid"},
		},
		"allowed version - passes version constraint": {
			exampleDir: "examples/allowedversion",
			config: &gomodguard.Configuration{
				Allowed: gomodguard.Allowed{
					"github.com/Masterminds/semver/v3": gomodguard.AllowedRule{
						Version: mustConstraint(t, ">= 3.0.0"),
					},
				},
			},
			wantEmpty: true,
		},
		"regex version - passes constraint": {
			exampleDir: "examples/regexversion",
			config: &gomodguard.Configuration{
				Blocked: gomodguard.Blocked{
					"golang\\.org/x/.*": gomodguard.BlockedRule{
						MatchType: gomodguard.RegexMatch,
						Version:   mustConstraint(t, "<= 0.15.0"),
						Reason:    "testing regex blocking with version constraint.",
					},
				},
			},
			wantEmpty: true,
		},
		"precedence - exact rule wins over overlapping regex rule": {
			exampleDir: "examples/alloptions",
			config: &gomodguard.Configuration{
				Blocked: gomodguard.Blocked{
					// Regex rule matches the same module but with a different reason
					"github\\.com/uudashr/.*": gomodguard.BlockedRule{
						MatchType: gomodguard.RegexMatch,
						Reason:    "regex catch-all should NOT be selected.",
					},
					// Exact rule should win due to higher precedence
					"github.com/uudashr/go-module": gomodguard.BlockedRule{
						MatchType:       gomodguard.ExactMatch,
						Recommendations: []string{"golang.org/x/mod"},
						Reason:          "exact rule should be selected.",
					},
				},
			},
			wantReasons: []string{
				"blocked_example.go:8:1 import of package `github.com/uudashr/go-module` is blocked because the " +
					"module is in the blocked modules list. `golang.org/x/mod` is a recommended module. " +
					"exact rule should be selected.",
			},
			notWantReasons: []string{
				"regex catch-all should NOT be selected",
			},
		},
		"precedence - prefix rule wins over overlapping regex rule": {
			exampleDir: "examples/alloptions",
			config: &gomodguard.Configuration{
				Blocked: gomodguard.Blocked{
					// Regex rule matches the same module
					"github\\.com/uudashr/.*": gomodguard.BlockedRule{
						MatchType: gomodguard.RegexMatch,
						Reason:    "regex catch-all should NOT be selected.",
					},
					// Prefix rule should win over regex
					"github.com/uudashr/": gomodguard.BlockedRule{
						MatchType: gomodguard.PrefixMatch,
						Reason:    "prefix rule should be selected.",
					},
				},
			},
			wantReasons: []string{
				"blocked_example.go:8:1 import of package `github.com/uudashr/go-module` is blocked because the " +
					"module is in the blocked modules list. prefix rule should be selected.",
			},
			notWantReasons: []string{
				"regex catch-all should NOT be selected",
			},
		},
		"precedence - longest prefix wins over shorter prefix": {
			exampleDir: "examples/alloptions",
			config: &gomodguard.Configuration{
				Blocked: gomodguard.Blocked{
					// Short prefix matches broadly
					"github.com/uudashr": gomodguard.BlockedRule{
						MatchType: gomodguard.PrefixMatch,
						Reason:    "short prefix should NOT be selected.",
					},
					// Longer prefix is more specific and should win
					"github.com/uudashr/go-module": gomodguard.BlockedRule{
						MatchType:       gomodguard.PrefixMatch,
						Recommendations: []string{"golang.org/x/mod"},
						Reason:          "longest prefix should be selected.",
					},
				},
			},
			wantReasons: []string{
				"blocked_example.go:8:1 import of package `github.com/uudashr/go-module` is blocked because the " +
					"module is in the blocked modules list. `golang.org/x/mod` is a recommended module. " +
					"longest prefix should be selected.",
			},
			notWantReasons: []string{
				"short prefix should NOT be selected",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Chdir(tt.exampleDir)

			reasons := processFiles(t, tt.config)

			if tt.wantEmpty {
				assert.Empty(t, reasons)
				return
			}

			for _, want := range tt.wantReasons {
				assert.Contains(t, reasons, want)
			}

			for _, notWant := range tt.notWantReasons {
				for _, r := range reasons {
					assert.NotContains(t, r, notWant)
				}
			}
		})
	}
}
