package cli_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v4"

	"github.com/ryancurrah/gomodguard/v2"
)

const examplesDir = "../../../../examples/"

func loadConfig(t *testing.T) *gomodguard.Configuration {
	t.Helper()

	data, err := os.ReadFile(".gomodguard.yaml")
	require.NoError(t, err)

	var config gomodguard.Configuration
	require.NoError(t, yaml.Unmarshal(data, &config))

	return &config
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

func TestProcessorYAMLConfig(t *testing.T) { //nolint:funlen
	tests := map[string]struct {
		exampleDir     string
		wantReasons    []string
		notWantReasons []string
		wantEmpty      bool
		wantParseErr   bool
		invalidYAML    string
	}{
		"all options - blocked by recommendation": {
			exampleDir: examplesDir + "alloptions",
			wantReasons: []string{
				"blocked_example.go:8:1 import of package `github.com/uudashr/go-module` is blocked because the " +
					"module is in the blocked modules list. `golang.org/x/mod` is a recommended module. `mod` " +
					"is the official go.mod parser library.",
				"blocked_example.go:7:1 import of package `github.com/mitchellh/go-homedir` is blocked because " +
					"the module is in the blocked modules list. version `v1.1.0` is blocked because it does not " +
					"meet the version constraint `<=1.1.0`. testing if blocked version constraint works.",
				"blocked_example.go:6:1 import of package `github.com/gofrs/uuid` is blocked because the " +
					"module is in the blocked modules list. `github.com/ryancurrah/gomodguard` is a recommended " +
					"module. testing if module is not blocked when it is recommended.",
			},
		},
		"allowed version - blocked by version constraint": {
			exampleDir: examplesDir + "allowedversion",
			wantReasons: []string{
				"example.go:3:1 import of package `github.com/Masterminds/semver/v3` is blocked because " +
					"version `v3.1.0` does not meet the allowed version constraint `>=3.2.0`.",
			},
		},
		"allowed version - invalid constraint is caught at parse time": {
			exampleDir:   examplesDir + "allowedversion",
			wantParseErr: true,
			invalidYAML:  "allowed:\n  - module: github.com/Masterminds/semver/v3\n    version: not-a-valid-constraint\n",
		},
		"empty allow list - allows all modules": {
			exampleDir: examplesDir + "emptyallowlist",
			wantEmpty:  true,
		},
		"indirect dependency - blocked": {
			exampleDir: examplesDir + "indirectdep",
			wantReasons: []string{
				"indirect_example.go:9:1 import of package `github.com/uudashr/go-module` is blocked because the " +
					"module is in the blocked modules list. `golang.org/x/mod` is a recommended module. `mod` " +
					"is the official go.mod parser library.",
				"indirect_example.go:6:1 import of package `github.com/gofrs/uuid` is blocked because the " +
					"module is in the blocked modules list. testing blocked indirect dependency.",
			},
		},
		"blocked invalid constraint is caught at parse time": {
			exampleDir:   examplesDir + "invalidconstraint",
			wantParseErr: true,
		},
		"regex - blocked": {
			exampleDir: examplesDir + "regextest",
			wantReasons: []string{
				"test.go:3:1 import of package `golang.org/x/mod/modfile` is blocked because the " +
					"module is in the blocked modules list. testing regex based blocking.",
			},
		},
		"regex version - blocked": {
			exampleDir: examplesDir + "regexversion",
			wantReasons: []string{
				"test.go:3:1 import of package `golang.org/x/mod/modfile` is blocked because the " +
					"module is in the blocked modules list. version `v0.16.0` is blocked because it does not " +
					"meet the version constraint `<=0.16.0`. testing regex blocking with version constraint.",
			},
		},
		"major version module is not blocked by base module rule": {
			exampleDir: examplesDir + "majorversion",
			wantReasons: []string{
				"example.go:4:1 import of package `github.com/gofrs/uuid` is blocked because the " +
					"module is in the blocked modules list. `github.com/gofrs/uuid/v5` is a recommended " +
					"module. testing that a major version module is not blocked by a rule targeting the base module.",
			},
			notWantReasons: []string{
				"import of package `github.com/gofrs/uuid/v5`",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Chdir(tt.exampleDir)

			if tt.wantParseErr {
				var yamlData []byte
				if tt.invalidYAML != "" {
					yamlData = []byte(tt.invalidYAML)
				} else {
					var err error

					yamlData, err = os.ReadFile(".gomodguard.yaml")
					require.NoError(t, err)
				}

				var cfg gomodguard.Configuration
				assert.Error(t, yaml.Unmarshal(yamlData, &cfg))

				return
			}

			config := loadConfig(t)

			reasons := processFiles(t, config)

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
