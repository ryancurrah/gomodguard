//nolint:scopelint
package gomodguard_test

import (
	"os"
	"strings"
	"testing"

	"github.com/ryancurrah/gomodguard"
	"github.com/ryancurrah/gomodguard/internal/filesearch"
)

func TestProcessorNewProcessor(t *testing.T) {
	_, err := gomodguard.NewProcessor(&gomodguard.Configuration{
		Allowed: gomodguard.Allowed{
			Modules: []string{
				"github.com/foo/bar",
			},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestProcessorProcessFiles(t *testing.T) { //nolint:funlen
	err := os.Chdir("_example/allOptions")
	if err != nil {
		t.Error(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	config := &gomodguard.Configuration{
		Allowed: gomodguard.Allowed{
			Modules: []string{
				"gopkg.in/yaml.v2",
				"github.com/go-xmlfmt/xmlfmt",
				"github.com/Masterminds/semver/v3",
				"github.com/ryancurrah/gomodguard",
			},
			Domains: []string{
				"golang.org",
			},
		},
		Blocked: gomodguard.Blocked{
			Modules: gomodguard.BlockedModules{
				{
					"github.com/uudashr/go-module": gomodguard.BlockedModule{
						Recommendations: []string{"golang.org/x/mod"},
						Reason:          "`mod` is the official go.mod parser library.",
					},
				},
				{
					"github.com/gofrs/uuid": gomodguard.BlockedModule{
						Recommendations: []string{"github.com/ryancurrah/gomodguard"},
						Reason:          "testing if module is not blocked when it is recommended.",
					},
				},
			},
			Versions: gomodguard.BlockedVersions{
				{
					"github.com/mitchellh/go-homedir": gomodguard.BlockedVersion{
						Version: "<= 1.1.0",
						Reason:  "testing if blocked version constraint works.",
					},
				},
			},
			LocalReplaceDirectives: true,
		},
	}

	processor, err := gomodguard.NewProcessor(config)
	if err != nil {
		t.Error(err)
	}

	filteredFiles := filesearch.Find(cwd, false, []string{"./..."})

	var tests = []struct {
		testName   string
		processor  gomodguard.Processor
		wantReason string
	}{
		{
			"module blocked because of recommendation",
			gomodguard.Processor{Config: config, Modfile: processor.Modfile},
			"blocked_example.go:9:1 import of package `github.com/uudashr/go-module` is blocked because the " +
				"module is in the blocked modules list. `golang.org/x/mod` is a recommended module. `mod` " +
				"is the official go.mod parser library.",
		},
		{
			"module blocked because of version constraint",
			gomodguard.Processor{Config: config, Modfile: processor.Modfile},
			"blocked_example.go:7:1 import of package `github.com/mitchellh/go-homedir` is blocked because " +
				"the module is in the blocked modules list. version `v1.1.0` is blocked because it does not " +
				"meet the version constraint `<= 1.1.0`. testing if blocked version constraint works.",
		},
		{
			"module blocked because of local replace directive",
			gomodguard.Processor{Config: config, Modfile: processor.Modfile},
			"blocked_example.go:8:1 import of package `github.com/ryancurrah/gomodguard` is blocked because " +
				"the module has a local replace directive.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			tt.processor.SetBlockedModules()

			results := tt.processor.ProcessFiles(filteredFiles)

			if len(results) == 0 {
				t.Fatal("result should be greater than zero")
			}

			foundWantReason := false
			allReasons := make([]string, 0, len(results))

			for _, result := range results {
				allReasons = append(allReasons, result.String())

				if strings.EqualFold(result.String(), tt.wantReason) {
					foundWantReason = true
				}
			}

			if !foundWantReason {
				t.Errorf("got '%+v' want '%s'", allReasons, tt.wantReason)
			}
		})
	}
}
