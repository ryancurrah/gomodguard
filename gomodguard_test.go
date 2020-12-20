// nolint:scopelint
package gomodguard_test

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/ryancurrah/gomodguard"
)

func TestBlockedModuleIsAllowed(t *testing.T) {
	var tests = []struct {
		testName          string
		blockedModule     gomodguard.BlockedModule
		currentModuleName string
		wantIsAllowed     bool
	}{
		{
			"blocked",
			gomodguard.BlockedModule{Recommendations: []string{"github.com/somerecommended/module"}},
			"github.com/ryancurrah/gomodguard",
			false,
		},
		{
			"allowed",
			gomodguard.BlockedModule{Recommendations: []string{"github.com/ryancurrah/gomodguard"}},
			"github.com/ryancurrah/gomodguard",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			isAllowed := tt.blockedModule.IsCurrentModuleARecommendation(tt.currentModuleName)
			if isAllowed != tt.wantIsAllowed {
				t.Errorf("got '%v' want '%v'", isAllowed, tt.wantIsAllowed)
			}
		})
	}
}

func TestBlockedModuleMessage(t *testing.T) {
	blockedWithNoRecommendation := "Some reason."
	blockedWithRecommendation := "`github.com/somerecommended/module` is a recommended module. Some reason."
	blockedWithRecommendations := "`github.com/somerecommended/module`, `github.com/someotherrecommended/module` and `github.com/someotherotherrecommended/module` are recommended modules. Some reason."

	var tests = []struct {
		testName          string
		blockedModule     gomodguard.BlockedModule
		currentModuleName string
		wantMessage       string
	}{
		{
			"blocked with no recommendation",
			gomodguard.BlockedModule{Recommendations: []string{}, Reason: "Some reason."},
			"github.com/ryancurrah/gomodguard",
			blockedWithNoRecommendation,
		},
		{
			"blocked with recommendation",
			gomodguard.BlockedModule{Recommendations: []string{"github.com/somerecommended/module"}, Reason: "Some reason."},
			"github.com/ryancurrah/gomodguard",
			blockedWithRecommendation,
		},
		{
			"blocked with multiple recommendations",
			gomodguard.BlockedModule{Recommendations: []string{"github.com/somerecommended/module", "github.com/someotherrecommended/module", "github.com/someotherotherrecommended/module"}, Reason: "Some reason."},
			"github.com/ryancurrah/gomodguard",
			blockedWithRecommendations,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			message := tt.blockedModule.Message()
			if !strings.EqualFold(message, tt.wantMessage) {
				t.Errorf("got '%s' want '%s'", message, tt.wantMessage)
			}
		})
	}
}

func TestBlockedModuleHasRecommendations(t *testing.T) {
	var tests = []struct {
		testName      string
		blockedModule gomodguard.BlockedModule
		wantIsAllowed bool
	}{
		{
			"does not have recommendations",
			gomodguard.BlockedModule{Recommendations: []string{}},
			false,
		},
		{
			"does have recommendations",
			gomodguard.BlockedModule{Recommendations: []string{"github.com/ryancurrah/gomodguard"}},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			hasRecommendations := tt.blockedModule.HasRecommendations()
			if hasRecommendations != tt.wantIsAllowed {
				t.Errorf("got '%v' want '%v'", hasRecommendations, tt.wantIsAllowed)
			}
		})
	}
}

func TestBlockedModulesGet(t *testing.T) {
	var tests = []struct {
		testName           string
		blockedModules     gomodguard.BlockedModules
		wantBlockedModules []string
	}{
		{
			"get all blocked module names",
			gomodguard.BlockedModules{{"github.com/someblocked/module": gomodguard.BlockedModule{Recommendations: []string{"github.com/ryancurrah/gomodguard"}}}},
			[]string{"github.com/someblocked/module"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			blockedModules := tt.blockedModules.Get()
			if !reflect.DeepEqual(blockedModules, tt.wantBlockedModules) {
				t.Errorf("got '%v' want '%v'", blockedModules, tt.wantBlockedModules)
			}
		})
	}
}

func TestBlockedVersionMessage(t *testing.T) {
	blockedWithVersionConstraint := "version `1.0.0` is blocked because it does not meet the version constraint `1.0.0`. Some reason."
	blockedWithVersionConstraintNoReason := "version `1.0.0` is blocked because it does not meet the version constraint `<= 1.0.0`."

	var tests = []struct {
		testName            string
		blockedVersion      gomodguard.BlockedVersion
		lintedModuleVersion string
		wantMessage         string
	}{
		{
			"blocked with version constraint",
			gomodguard.BlockedVersion{Version: "1.0.0", Reason: "Some reason."},
			"1.0.0",
			blockedWithVersionConstraint,
		},
		{
			"blocked with version constraint and no reason",
			gomodguard.BlockedVersion{Version: "<= 1.0.0"},
			"1.0.0",
			blockedWithVersionConstraintNoReason,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			message := tt.blockedVersion.Message(tt.lintedModuleVersion)
			if !strings.EqualFold(message, tt.wantMessage) {
				t.Errorf("got '%s' want '%s'", message, tt.wantMessage)
			}
		})
	}
}

func TestBlockedModulesGetBlockedModule(t *testing.T) {
	var tests = []struct {
		testName          string
		blockedModules    gomodguard.BlockedModules
		currentModuleName string
		lintedModuleName  string
		wantIsAllowed     bool
	}{
		{
			"blocked",
			gomodguard.BlockedModules{{"github.com/someblocked/module": gomodguard.BlockedModule{Recommendations: []string{"github.com/someother/module"}}}},
			"github.com/ryancurrah/gomodguard",
			"github.com/someblocked/module",
			false,
		},
		{
			"allowed",
			gomodguard.BlockedModules{{"github.com/someblocked/module": gomodguard.BlockedModule{Recommendations: []string{"github.com/ryancurrah/gomodguard"}}}},
			"github.com/ryancurrah/gomodguard",
			"github.com/someblocked/module",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			blockedModule := tt.blockedModules.GetBlockReason(tt.lintedModuleName)
			if blockedModule.IsCurrentModuleARecommendation(tt.currentModuleName) != tt.wantIsAllowed {
				t.Errorf("got '%+v' want '%+v'", blockedModule.IsCurrentModuleARecommendation(tt.currentModuleName), tt.wantIsAllowed)
			}
		})
	}
}

func TestAllowedIsAllowedModule(t *testing.T) {
	var tests = []struct {
		testName            string
		allowedModules      gomodguard.Allowed
		lintedModuleName    string
		wantIsAllowedModule bool
	}{
		{
			"module is allowed",
			gomodguard.Allowed{Modules: []string{"github.com/someallowed/module"}},
			"github.com/someallowed/module",
			true,
		},
		{
			"module not allowed",
			gomodguard.Allowed{},
			"github.com/someblocked/module",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			isAllowedModule := tt.allowedModules.IsAllowedModule(tt.lintedModuleName)
			if !reflect.DeepEqual(isAllowedModule, tt.wantIsAllowedModule) {
				t.Errorf("got '%v' want '%v'", isAllowedModule, tt.wantIsAllowedModule)
			}
		})
	}
}

func TestAllowedIsAllowedModuleDomain(t *testing.T) {
	var tests = []struct {
		testName                  string
		allowedModules            gomodguard.Allowed
		lintedModuleName          string
		wantIsAllowedModuleDomain bool
	}{
		{
			"module is allowed",
			gomodguard.Allowed{Domains: []string{"github.com"}},
			"github.com/someallowed/module",
			true,
		},
		{
			"module not allowed",
			gomodguard.Allowed{},
			"github.com/someblocked/module",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			isAllowedModuleDomain := tt.allowedModules.IsAllowedModuleDomain(tt.lintedModuleName)
			if !reflect.DeepEqual(isAllowedModuleDomain, tt.wantIsAllowedModuleDomain) {
				t.Errorf("got '%v' want '%v'", isAllowedModuleDomain, tt.wantIsAllowedModuleDomain)
			}
		})
	}
}

func TestResultString(t *testing.T) {
	var tests = []struct {
		testName   string
		result     gomodguard.Result
		wantString string
	}{
		{
			"reason lint failed",
			gomodguard.Result{FileName: "test.go", LineNumber: 1, Reason: "Some reason."},
			"test.go:1:1 Some reason.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			result := tt.result.String()
			if !strings.EqualFold(result, tt.wantString) {
				t.Errorf("got '%s' want '%s'", result, tt.wantString)
			}
		})
	}
}

func TestProcessorNewProcessor(t *testing.T) {
	config, _, err := setup()
	if err != nil {
		t.Error(err)
	}

	_, err = gomodguard.NewProcessor(config)
	if err != nil {
		t.Error(err)
	}
}

func TestProcessorProcessFiles(t *testing.T) {
	config, cwd, err := setup()
	if err != nil {
		t.Error(err)
	}

	processor, err := gomodguard.NewProcessor(config)
	if err != nil {
		t.Error(err)
	}

	filteredFiles := gomodguard.GetFilteredFiles(cwd, false, []string{"./..."})

	var tests = []struct {
		testName   string
		processor  gomodguard.Processor
		wantReason string
	}{
		{
			"module blocked because of recommendation",
			gomodguard.Processor{Config: config, Modfile: processor.Modfile, Result: []gomodguard.Result{}},
			"blocked_example.go:9:1 import of package `github.com/uudashr/go-module` is blocked because the module is in the blocked modules list. `golang.org/x/mod` is a recommended module. `mod` is the official go.mod parser library.",
		},
		{
			"module blocked because of version constraint",
			gomodguard.Processor{Config: config, Modfile: processor.Modfile, Result: []gomodguard.Result{}},
			"blocked_example.go:7:1 import of package `github.com/mitchellh/go-homedir` is blocked because the module is in the blocked modules list. version `v1.1.0` is blocked because it does not meet the version constraint `<= 1.1.0`. testing if blocked version constraint works.",
		},
		{
			"module blocked because of local replace directive",
			gomodguard.Processor{Config: config, Modfile: processor.Modfile, Result: []gomodguard.Result{}},
			"blocked_example.go:8:1 import of package `github.com/ryancurrah/gomodguard` is blocked because the module has a local replace directive.",
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

func setup() (*gomodguard.Configuration, string, error) {
	os.Chdir("_example")
	cwd, _ := os.Getwd()

	config, err := gomodguard.GetConfig(".gomodguard.yaml")
	if err != nil {
		return nil, "", err
	}

	return config, cwd, nil
}
