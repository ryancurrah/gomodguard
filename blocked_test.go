package gomodguard_test

import (
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
			gomodguard.BlockedModule{
				Recommendations: []string{
					"github.com/somerecommended/module",
				},
			},
			"github.com/ryancurrah/gomodguard",
			false,
		},
		{
			"allowed",
			gomodguard.BlockedModule{
				Recommendations: []string{
					"github.com/ryancurrah/gomodguard",
				},
			},
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
	blockedWithRecommendations := "`github.com/somerecommended/module`, `github.com/someotherrecommended/module` " +
		"and `github.com/someotherotherrecommended/module` are recommended modules. Some reason."

	var tests = []struct {
		testName          string
		blockedModule     gomodguard.BlockedModule
		currentModuleName string
		wantMessage       string
	}{
		{
			"blocked with no recommendation",
			gomodguard.BlockedModule{
				Recommendations: []string{},
				Reason:          "Some reason.",
			},
			"github.com/ryancurrah/gomodguard",
			blockedWithNoRecommendation,
		},
		{
			"blocked with recommendation",
			gomodguard.BlockedModule{
				Recommendations: []string{
					"github.com/somerecommended/module",
				},
				Reason: "Some reason.",
			},
			"github.com/ryancurrah/gomodguard",
			blockedWithRecommendation,
		},
		{
			"blocked with multiple recommendations",
			gomodguard.BlockedModule{
				Recommendations: []string{
					"github.com/somerecommended/module",
					"github.com/someotherrecommended/module",
					"github.com/someotherotherrecommended/module",
				},
				Reason: "Some reason.",
			},
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
			gomodguard.BlockedModule{
				Recommendations: []string{
					"github.com/ryancurrah/gomodguard",
				},
			},
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
			gomodguard.BlockedModules{
				{
					"github.com/someblocked/module": gomodguard.BlockedModule{
						Recommendations: []string{
							"github.com/ryancurrah/gomodguard",
						},
					},
				},
			},
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
	blockedWithVersionConstraint := "version `1.0.0` is blocked because it does not meet the version constraint " +
		"`1.0.0`. Some reason."
	blockedWithVersionConstraintNoReason := "version `1.0.0` is blocked because it does not meet the version " +
		"constraint `<= 1.0.0`."

	var tests = []struct {
		testName            string
		blockedVersion      gomodguard.BlockedVersion
		lintedModuleVersion string
		wantMessage         string
	}{
		{
			"blocked with version constraint",
			gomodguard.BlockedVersion{
				Version: "1.0.0",
				Reason:  "Some reason.",
			},
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
			gomodguard.BlockedModules{
				{
					"github.com/someblocked/module": gomodguard.BlockedModule{
						Recommendations: []string{
							"github.com/someother/module",
						},
					},
				},
			},
			"github.com/ryancurrah/gomodguard",
			"github.com/someblocked/module",
			false,
		},
		{
			"allowed",
			gomodguard.BlockedModules{
				{
					"github.com/someblocked/module": gomodguard.BlockedModule{
						Recommendations: []string{
							"github.com/ryancurrah/gomodguard",
						},
					},
				},
			},
			"github.com/ryancurrah/gomodguard",
			"github.com/someblocked/module",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			blockedModule := tt.blockedModules.GetBlockReason(tt.lintedModuleName)
			if blockedModule.IsCurrentModuleARecommendation(tt.currentModuleName) != tt.wantIsAllowed {
				t.Errorf("got '%+v' want '%+v'", blockedModule.IsCurrentModuleARecommendation(tt.currentModuleName),
					tt.wantIsAllowed)
			}
		})
	}
}
